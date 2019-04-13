package protocol

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/odysseyhack/planet-society/protocol/models"
	log "github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/parser"
)

const (
	connectionChannelSize = 16
)

type Protocol struct {
	quit             chan struct{}
	Connections      chan Conn
	authorization    AuthorizationPlugin
	TransactionQueue *Queue
	plugins          Plugins
}

func NewProtocol(authorization AuthorizationPlugin) *Protocol {
	return &Protocol{
		authorization:    authorization,
		quit:             make(chan struct{}),
		Connections:      make(chan Conn, connectionChannelSize),
		TransactionQueue: NewQueue(),
	}
}

func (p *Protocol) Stop() {
	p.quit <- struct{}{}
}

func (p *Protocol) Loop() {
	log.Debugln("protocol: starting event loop")
	for {
		select {
		case <-p.quit:
			return
		case conn := <-p.Connections:
			log.Debugln("protocol: handling new connection")
			go p.handleConn(conn)
		}
	}
}

func (p *Protocol) handleConn(c Conn) {
	log.Debugln("protocol handling new connection")
	defer func() {
		if err := c.Close(); err != nil {
			log.Warningln("protocol: failed to close connection:", err)
		}
	}()

	for {
		msg, err := c.Read()
		if err != nil {
			break
		}
		p.handleMessage(c, msg)
	}
}

func (p *Protocol) handleMessage(c Conn, msg *Message) {
	switch msg.Header.Topic {
	case TopicPreTransactionRequest:
		p.handlePreTransactionRequest(c, msg)
	case TopicTransactionRequest:
		p.handleTransactionRequest(c, msg)
	}
}

func (p *Protocol) handleTransactionRequest(c Conn, msg *Message) {
	var transactionRequest models.TransactionRequest
	if err := gob.NewDecoder(bytes.NewBuffer(msg.Body.Payload)).Decode(&transactionRequest); err != nil {
		errMsg := "decoding payload failed"
		sendTransactionReply(c, msg, &models.TransactionReply{Error: &errMsg})
		log.Warningln("transaction request: invalid payload:", err)
		return
	}

	entry, ok := p.TransactionQueue.Get(transactionRequest.TransactionID.Key)
	if !ok {
		log.Warningln("transaction is not in TransactionQueue id:", transactionRequest.TransactionID)
		errMsg := "transaction does not exist in queue"
		sendTransactionReply(c, msg, &models.TransactionReply{Error: &errMsg})
		return
	}

	data, err := parseQuery(transactionRequest.Query)
	if err != nil {
		log.Warningln("transaction failed to parse query id:", transactionRequest.TransactionID)
		errMsg := "query parsing failed"
		sendTransactionReply(c, msg, &models.TransactionReply{Error: &errMsg})
		return
	}

	authData := generateNotificationRequest(&transactionRequest, data, entry)

	log.Infoln("calling authorization plugin for id:", transactionRequest.TransactionID.Key.String())
	authReply, err := p.authorization.Authorize(authData)
	if err != nil {
		log.Warningf("transaction failed to authorize id=%q , err=%q", transactionRequest.TransactionID.Key.String(), err)
		errMsg := "not authorized"
		sendTransactionReply(c, msg, &models.TransactionReply{Error: &errMsg})
		return
	}

	if !authReply.Accepted {
		log.Warningf("transaction was not authorized id=%q", transactionRequest.TransactionID.Key.String())
		errMsg := "not authorized"
		sendTransactionReply(c, msg, &models.TransactionReply{Error: &errMsg})
		return
	}

	content, err := post(transactionRequest.Query, &transactionRequest, entry)
	if err != nil {
		log.Warningf("transaction failed to post transaction id=%q, err=%q", transactionRequest.TransactionID.Key.String(), err)
		errMsg := "transaction commitment failed"
		sendTransactionReply(c, msg, &models.TransactionReply{Error: &errMsg})
		return
	}
	sendTransactionReply(c, msg, &models.TransactionReply{Content: &content})
}

func generateNotificationRequest(request *models.TransactionRequest, c []CollectionData, e *Entry) *models.PermissionNotificationRequest {
	ret := &models.PermissionNotificationRequest{
		RequesterName:      e.RequesterName,
		Date:               time.Now().Format(time.RFC3339),
		Description:        request.Description,
		Title:              request.Title,
		RequesterPublicKey: e.RequesterPublicKey.String(),
		TransactionID:      e.TransactionID.String(),
		Analysis:           []string{"personal data is GDPR protected data", "banking details is sensitive data"},
		Verification:       []string{"digid.nl", "planet-blockchain", "kvk"},
	}

	for i := range c {
		ret.Item = append(ret.Item, models.ItemField{Item: c[i].Structure, Fields: c[i].Fields})
	}

	return ret
}

func (p *Protocol) handlePreTransactionRequest(c Conn, msg *Message) {
	var preTransactionRequest models.PreTransactionRequest
	if err := gob.NewDecoder(bytes.NewBuffer(msg.Body.Payload)).Decode(&preTransactionRequest); err != nil {
		sendPreTransactionReply(c, msg, false)
		log.Warningln("pre transaction request: invalid payload:", err)
		return
	}

	log.Infoln("protocol: validating pre transaction request with plugins")
	if !p.plugins.ValidatePreTransaction(&preTransactionRequest) {
		log.Warningln("protocol: request didn't pass validation")
	}

	entry := &Entry{
		TransactionID:      preTransactionRequest.TransactionID.Key,
		RequesterName:      preTransactionRequest.Requester,
		RequesterPublicKey: preTransactionRequest.MainPublicKey.Key,
	}

	if err := p.TransactionQueue.Add(entry); err != nil {
		sendPreTransactionReply(c, msg, false)
		log.Warningln("protocol: adding to TransactionQueue failed:", err)
	}
	log.Infoln("protocol: added new transaction to TransactionQueue")
	sendPreTransactionReply(c, msg, true)
}

func sendPreTransactionReply(c Conn, msg *Message, ok bool) {
	var buffer bytes.Buffer

	reply := &models.PreTransactionReply{
		Success: ok,
	}

	if err := gob.NewEncoder(&buffer).Encode(reply); err != nil {
		return
	}

	_ = c.Write(&Message{
		Header: Header{
			Topic:       TopicPreTransactionReply,
			Destination: msg.Header.Source,
			Source:      msg.Header.Destination,
		},
		Body: Body{
			Payload: buffer.Bytes(),
		},
	})
}

func sendTransactionReply(c Conn, msg *Message, reply *models.TransactionReply) {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(reply); err != nil {
		return
	}

	_ = c.Write(&Message{
		Header: Header{
			Topic:       TopicTransactionReply,
			Destination: msg.Header.Source,
			Source:      msg.Header.Destination,
		},
		Body: Body{
			Payload: buffer.Bytes(),
		},
	})
}

type CollectionData struct {
	Structure string
	Fields    []string
}

func parseQuery(query string) (c []CollectionData, err error) {
	doc, gqlErr := parser.ParseQuery(&ast.Source{Input: query})
	if gqlErr != nil {
		return c, gqlErr
	}
	c = parse(doc.Operations)
	return c, err
}

func parse(operations ast.OperationList) (c []CollectionData) {
	for i := range operations {
		for _, op := range operations[i].SelectionSet {
			if fragment, ok := op.(*ast.Field); ok {
				data := &CollectionData{Structure: fragment.Name}
				getFields(data, fragment.SelectionSet)
				c = append(c, *data)
			}
		}
	}
	return c
}

func getFields(c *CollectionData, s ast.SelectionSet) {
	for _, op := range s {
		if field, ok := op.(*ast.Field); ok {
			c.Fields = append(c.Fields, field.Name)
		}
	}
}

type Request struct {
	Query string `json:"query"`
}

func post(query string, transactionRequest *models.TransactionRequest, entry *Entry) (string, error) {
	r := Request{
		Query: query,
	}

	requestBody, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("encode: %s", err.Error())
	}

	client := http.Client{}
	request, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/query", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	request.Header.Add("title", transactionRequest.Title)
	request.Header.Add("description", transactionRequest.Description)
	request.Header.Add("TransactionID", transactionRequest.TransactionID.Key.String())
	request.Header.Add("signature", transactionRequest.Signature)
	request.Header.Add("requester", entry.RequesterPublicKey.String())
	request.Header.Add("requester-name", entry.RequesterName)
	request.Header.Set("Content-Type", "application/json")
	rawResponse, err := client.Do(request)
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(rawResponse.Body)
	if err != nil {
		return "", err
	}

	var reply Reply
	if err := json.Unmarshal(content, &reply); err != nil {
		log.Warningln("failed to unmarshal reply from permission:", err)
	}

	if len(reply.Errors) > 0 {
		return "", fmt.Errorf("query failed: %v", reply.Errors)
	}

	return dataReplyTorString(reply.Data)
}

type ErrorsReply struct {
	Message string `json:"message,omitempty"`
}

type Reply struct {
	Errors []ErrorsReply                `json:"errors,omitempty"`
	Data   map[string]map[string]string `json:"data,omitempty"`
}

func dataReplyTorString(data map[string]map[string]string) (string, error) {
	ret, err := json.Marshal(data)
	return string(ret), err
}
