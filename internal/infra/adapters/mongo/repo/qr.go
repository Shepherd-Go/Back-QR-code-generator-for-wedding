package repo

import (
	"context"
	"fmt"

	"github.com/andresxlp/qr-system/internal/domain/dto"
	"github.com/andresxlp/qr-system/internal/domain/ports/repo"
	"github.com/andresxlp/qr-system/internal/infra/adapters/mongo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type qr struct {
	dbClient models.DBClientWrite
}

func NewQr(dbClient models.DBClientWrite) repo.QR {
	return &qr{dbClient}
}

func (q *qr) Create(ctx context.Context, qr models.Guest) (string, error) {
	db := q.dbClient.Collection("qr-codes")
	objectID, err := db.InsertOne(ctx, qr)
	if err != nil {
		return "", err
	}

	id := objectID.InsertedID.(primitive.ObjectID).Hex()

	return id, err
}

func (q *qr) ValidateQrCode(ctx context.Context, id primitive.ObjectID) (dto.Guest, error) {
	db := q.dbClient.Collection("qr-codes")
	filter := bson.D{{"_id", id}}

	infoGuest := models.Guest{}
	if err := db.FindOne(ctx, filter).Decode(&infoGuest); err != nil {
		if err == mongo.ErrNoDocuments {
			return dto.Guest{}, fmt.Errorf("this qr-code not exist")
		}
	}

	return infoGuest.ToDomainDTO(), nil
}

func (q *qr) ConfirmInvitation(ctx context.Context, id primitive.ObjectID) error {
	db := q.dbClient.Collection("qr-codes")

	_, err := db.UpdateOne(ctx, bson.D{{"_id", id}}, bson.D{{"$set", bson.D{{"status", "Used"}}}})
	if err != nil {
		return err
	}

	return nil
}

/*func (q qr) CountQRCodeUsed(ctx context.Context, emailOwner string) (int64, error) {
	db := q.dbClient.Collection("qr")
	filter := bson.D{{"created_by", emailOwner}, {"status", "Used"}}
	totalQRCodeUsed, err := db.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return totalQRCodeUsed, nil
}*/
