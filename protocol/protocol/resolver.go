package protocol

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/database"
	"github.com/odysseyhack/planet-society/protocol/models"
)

type Resolver struct {
	db *database.Database
}

func NewResolver(db *database.Database) *Resolver {
	return &Resolver{
		db: db,
	}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) PersonalDetailsUpdate(ctx context.Context, input models.PersonalDetailsInput) (*models.PersonalDetails, error) {
	details, err := r.db.PersonalDetailsUpdate(input)
	return &details, err
}

func (r *mutationResolver) ContactAdd(ctx context.Context, contacts models.ContactInput) (*models.Contact, error) {
	newContact, err := r.db.ContactAdd(contacts)
	return &newContact, err
}

func (r *mutationResolver) ContactDel(ctx context.Context, id string) (string, error) {
	return r.db.ContactDel(id)
}

func (r *mutationResolver) AddressAdd(ctx context.Context, addresses models.AddressInput) (*models.Address, error) {
	newAddress, err := r.db.AddressAdd(addresses)
	return &newAddress, err
}

func (r *mutationResolver) AddressDel(ctx context.Context, id string) (string, error) {
	return r.db.AddressDel(id)
}

func (r *mutationResolver) IdentityAdd(ctx context.Context, identities models.IdentityInput) (*models.Identity, error) {
	panic("not implemented")
}

func (r *mutationResolver) IdentityDel(ctx context.Context, id string) (string, error) {
	return r.db.IdentityDel(id)
}

func (r *mutationResolver) PaymentCardAdd(ctx context.Context, paymentCards models.PaymentCardInput) (*models.PaymentCard, error) {
	panic("not implemented")
}
func (r *mutationResolver) PaymentCardDel(ctx context.Context, id string) (string, error) {
	return r.db.PaymentCardDel(id)
}
func (r *mutationResolver) PassportAdd(ctx context.Context, passports models.PassportInput) (*models.Passport, error) {
	panic("not implemented")
}
func (r *mutationResolver) PassportDel(ctx context.Context, id string) (string, error) {
	return r.db.PassportDel(id)
}
func (r *mutationResolver) IdentityDocumentAdd(ctx context.Context, identityDocument models.IdentityDocumentInput) (*models.IdentityDocument, error) {
	panic("not implemented")
}
func (r *mutationResolver) IdentityDocumentDel(ctx context.Context, id string) (string, error) {
	return r.db.IdentityDocumentDel(id)
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) PersonalDetails(ctx context.Context) (*models.PersonalDetails, error) {
	t, err := transact(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return &t.PersonalDetails, nil
}
func (r *queryResolver) Address(ctx context.Context) (*models.Address, error) {
	t, err := transact(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return &t.Address, nil
}

func (r *queryResolver) PaymentCard(ctx context.Context) (*models.PaymentCard, error) {
	t, err := transact(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return &t.PaymentCard, nil
}
func (r *queryResolver) Passport(ctx context.Context) (*models.Passport, error) {
	t, err := transact(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return &t.Passport, nil
}
func (r *queryResolver) IdentityDocument(ctx context.Context) (*models.IdentityDocument, error) {
	t, err := transact(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return &t.IdentityDocument, nil
}
func (r *queryResolver) Identity(ctx context.Context) ([]models.Identity, error) {
	return r.db.IdentityList()
}

func (r *queryResolver) PermissionListByPublicKey(ctx context.Context, publicKey models.Key32) (ret []models.Permission, err error) {
	list, err := r.db.PermissionList()
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].RequesterPublicKey.Key.Equal(publicKey.Key) {
			ret = append(ret, list[i])
		}
	}
	return list, err
}

func (r *queryResolver) PermissionListByResource(ctx context.Context, id string) (ret []models.Permission, err error) {
	list, err := r.db.PermissionList()
	if err != nil {
		return nil, err
	}
	for i := range list {
		for j := range list[i].PermissionNodes {
			if list[i].PermissionNodes[j].NodeID == id {
				ret = append(ret, list[i])
				break
			}
		}
	}
	return list, err
}

func (r *queryResolver) PermissionList(ctx context.Context) ([]models.Permission, error) {
	a, err := r.db.PermissionList()
	return a, err
}

func (r *queryResolver) PaymentCardList(ctx context.Context, identity string) ([]models.PaymentCard, error) {
	return r.db.PaymentCardList(identity)
}
func (r *queryResolver) PassportList(ctx context.Context, identity string) ([]models.Passport, error) {
	return r.db.PassportList(identity)
}
func (r *queryResolver) IdentityDocumentList(ctx context.Context, identity string) ([]models.IdentityDocument, error) {
	return r.db.IdentityDocumentList(identity)
}

// dummy
type Transaction struct {
	Address          models.Address
	PaymentCard      models.PaymentCard
	Passport         models.Passport
	IdentityDocument models.IdentityDocument
	PersonalDetails  models.PersonalDetails
}

var (
	cache  = make(map[string]*Transaction)
	locker sync.Mutex
)

func randomTransaction() *models.Permission {
	tID := cryptography.RandomKey32()
	requesterPublicKey := cryptography.RandomKey32()
	RequesterSignatureKey := cryptography.RandomKey32()
	ResponderSignature := cryptography.RandomKey32()
	RequesterSignature := cryptography.RandomKey32()

	return &models.Permission{
		TransactionID:         tID.String(),
		Expiration:            time.Now().Add(time.Hour * 200).Format(time.RFC3339),
		RequesterPublicKey:    models.Key32{Key: requesterPublicKey},
		RequesterSignatureKey: models.Key32{Key: RequesterSignatureKey},
		Reason:                "please provide data to finalize shipping",
		ResponderSignature:    ResponderSignature.String(),
		RequesterSignature:    RequesterSignature.String(),
		Revokable:             false,
		LawApplying:           "European Union",
	}
}

func transact(ctx context.Context, db *database.Database) (*Transaction, error) {
	locker.Lock()
	defer locker.Unlock()

	transactionID, ok := ctx.Value("transaction-key").(string)
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	if t, ok := cache[transactionID]; ok {
		return t, nil
	}

	transaction, err := fillTransaction(db)
	if err != nil {
		return nil, err
	}

	cache[transactionID] = transaction
	tr := randomTransaction()
	fields := graphql.CollectAllFields(ctx)
	for _, field := range fields {
		switch field {
		case "personalDetails":
			tr.PermissionNodes = append(tr.PermissionNodes, models.PermissionNodes{NodeID: transaction.PersonalDetails.ID})
		case "address":
			tr.PermissionNodes = append(tr.PermissionNodes, models.PermissionNodes{NodeID: transaction.Address.ID})
		case "paymentCard":
			tr.PermissionNodes = append(tr.PermissionNodes, models.PermissionNodes{NodeID: transaction.PaymentCard.ID})
		case "passport":
			tr.PermissionNodes = append(tr.PermissionNodes, models.PermissionNodes{NodeID: transaction.Passport.ID})
		case "identityDocument":
			tr.PermissionNodes = append(tr.PermissionNodes, models.PermissionNodes{NodeID: transaction.IdentityDocument.ID})
		}
	}
	if _, err := db.PermissionAdd(*tr); err != nil {
		return nil, err
	}
	return transaction, nil
}

func fillTransaction(db *database.Database) (transaction *Transaction, err error) {
	transaction = &Transaction{}
	randS := rand.New(rand.NewSource(time.Now().UnixNano()))

	if transaction.PersonalDetails, err = db.PersonalDetails(); err != nil {
		return nil, err
	}

	identities, err := db.IdentityList()
	if err != nil {
		return nil, err
	}
	identityID := identities[randS.Uint32()%uint32(len(identities))].ID

	as, err := db.AddressList(identityID)
	if err != nil {
		return nil, err
	}
	number4 := randS.Uint32() % uint32(len(as))
	transaction.Address = as[number4]

	if err := fillTransactionWithDocuments(db, transaction); err != nil {
		return nil, err
	}

	return transaction, nil
}

func fillTransactionWithDocuments(db *database.Database, transaction *Transaction) error {
	randS := rand.New(rand.NewSource(time.Now().UnixNano()))

	identities, err := db.IdentityList()
	if err != nil {
		return err
	}
	identityID := identities[randS.Uint32()%uint32(len(identities))].ID

	ids, err := db.IdentityDocumentList(identityID)
	if err != nil {
		return err
	}
	number1 := randS.Uint32() % uint32(len(ids))
	transaction.IdentityDocument = ids[number1]

	ps, err := db.PassportList(identityID)
	if err != nil {
		return err
	}
	number2 := randS.Uint32() % uint32(len(ps))
	transaction.Passport = ps[number2]
	pc, err := db.PaymentCardList(identityID)
	if err != nil {
		return err
	}
	number3 := randS.Uint32() % uint32(len(pc))
	transaction.PaymentCard = pc[number3]
	return nil
}
