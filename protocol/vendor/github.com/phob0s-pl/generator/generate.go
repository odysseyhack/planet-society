package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/odysseyhack/planet-society/protocol/cryptography"
	"github.com/odysseyhack/planet-society/protocol/database"
	"github.com/odysseyhack/planet-society/protocol/models"
)

type Generator struct {
	rand       *rand.Rand
	identities []models.Identity
}

func NewGenerator() *Generator {
	return &Generator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *Generator) identityCreate() {

}

func (g *Generator) Generate(db *database.Database) error {
	if err := g.generatePersonalDetails(db); err != nil {
		return err
	}

	if err := g.generateIdentities(db); err != nil {
		return err
	}

	if err := g.generateContacts(db); err != nil {
		return err
	}

	if err := g.generateAddresses(db); err != nil {
		return err
	}

	if err := g.generateDocuments(db); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generatePersonalDetails(db *database.Database) error {
	name := g.randomName()
	surname := g.randomSurname()
	date := time.Now().Format(time.RFC3339)
	country := g.randomCountry()

	pd := models.PersonalDetailsInput{
		Name:      &name,
		Surname:   &surname,
		BirthDate: &date,
		Country:   &country,
	}

	_, err := db.PersonalDetailsUpdate(pd)
	return err
}

func (g *Generator) generateIdentities(db *database.Database) error {
	var newIdentities []models.IdentityInput
	identityNames := gidentityNames()
	for i := range identityNames {
		newIdentity := models.IdentityInput{
			DisplayName: identityNames[i],
		}
		newIdentities = append(newIdentities, newIdentity)
		identity, err := db.IdentityAdd(newIdentity)
		if err != nil {
			return err
		}
		g.identities = append(g.identities, identity)

	}
	return nil
}

func (g *Generator) generateContacts(db *database.Database) error {
	for i := range g.identities {
		numberOfContacts := (g.rand.Uint32() % 5) + 1

		for x := uint32(0); x < numberOfContacts; x++ {
			name := g.randomName()
			surname := g.randomSurname()
			display := fmt.Sprintf("%s.%s", name, surname)
			country := g.randomCountry()
			newContact := models.ContactInput{Identity: g.identities[i].ID, PublicKey: models.Key32{Key: cryptography.RandomKey32()}, Name: &name, DisplayName: display, Country: &country, Surname: &surname}
			if _, err := db.ContactAdd(newContact); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) generateAddresses(db *database.Database) error {
	for i := range g.identities {
		numberOfAddresses := (g.rand.Uint32() % 2) + 1
		for x := uint32(0); x < numberOfAddresses; x++ {
			country := g.randomCountry()
			city := g.randomCity()
			street := g.randomStreet()
			display := fmt.Sprintf("%s.%s", country, city)
			newAddress := models.AddressInput{Identity: g.identities[i].ID, City: &city, DisplayName: display, Country: &country, Street: &street}
			if _, err := db.AddressAdd(newAddress); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) generatePaymentCard(db *database.Database) error {
	for i := range g.identities {
		name := g.randomName()
		surname := g.randomSurname()
		paymentCard := models.PaymentCardInput{Identity: g.identities[i].ID, Name: name, Surname: surname, Number: fmt.Sprintf("%d", g.rand.Uint32()), DisplayName: fmt.Sprintf("my_card"), Expiration: time.Now().Format(time.RFC3339), SecurityCode: fmt.Sprintf("%d", g.rand.Uint32()%1000), Currency: g.randomCurrency()}
		if _, err := db.PaymentCardAdd(paymentCard); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) generatePassport(db *database.Database) error {
	for i := range g.identities {
		name := g.randomName()
		surname := g.randomSurname()
		passport := models.PassportInput{Identity: g.identities[i].ID, Name: name, Surname: surname, Number: fmt.Sprintf("%d", g.rand.Uint32()), Expiration: time.Now().Format(time.RFC3339), DisplayName: "my_passport", Country: g.randomCountry()}
		if _, err := db.PassportAdd(passport); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) generateIdentityDocument(db *database.Database) error {
	for i := range g.identities {
		name := g.randomName()
		surname := g.randomSurname()

		identityCard := models.IdentityDocumentInput{Identity: g.identities[i].ID, Name: name, Surname: surname, Number: fmt.Sprintf("%d", g.rand.Uint32()), Country: g.randomCountry(), DisplayName: "my_id_document", Expiration: time.Now().Format(time.RFC3339)}
		if _, err := db.IdentityDocumentAdd(identityCard); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) generateDocuments(db *database.Database) error {
	if err := g.generateIdentityDocument(db); err != nil {
		return err
	}
	if err := g.generatePassport(db); err != nil {
		return err
	}
	if err := g.generatePaymentCard(db); err != nil {
		return err
	}
	return nil
}

func (g *Generator) randomName() string {
	names := gnames()
	pos := g.rand.Uint32() % uint32(len(names))
	return names[pos]
}

func (g *Generator) randomSurname() string {
	surnames := gsurnames()
	pos := g.rand.Uint32() % uint32(len(surnames))
	return surnames[pos]
}

func (g *Generator) randomStreet() string {
	streets := gstreets()
	pos := g.rand.Uint32() % uint32(len(streets))
	return streets[pos]
}

func (g *Generator) randomCity() string {
	cities := gcities()
	pos := g.rand.Uint32() % uint32(len(cities))
	return cities[pos]
}

func (g *Generator) randomCountry() string {
	countries := gcountries()
	pos := g.rand.Uint32() % uint32(len(countries))
	return countries[pos]
}

func (g *Generator) randomCurrency() string {
	currencies := gcountries()
	pos := g.rand.Uint32() % uint32(len(currencies))
	return currencies[pos]
}
