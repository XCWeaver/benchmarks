// ! **********************************************************************************************************************************************
// ! **                                                                                                                                          **
// ! **          This service only has a queue-consumer functionality, but we found it to be disabled in the original implementation.            **
// ! **          However, we wanted to have this implemented as an example of a functioning queue-subsystem between two services.                **
// ! **                                                                                                                                          **
// ! **********************************************************************************************************************************************
package services

import (
	"encoding/json"
	"log"

	"trainticket/pkg/util"

	"github.com/ServiceWeaver/weaver"
	"gitlab.mpi-sws.org/cld/blueprint/blueprint-compiler/stdlib/components"
)

type DeliveryService interface {
	Entry()
}

type deliveryService struct {
	weaver.Implements[DeliveryService]
	deliveryQueue components.Queue
	//Mongo
	deliveryDB components.NoSQLDatabase
}

func (d *deliveryService) ProcessDelivery(payload []byte) {
	var delivery util.Delivery
	err := json.Unmarshal(payload, &delivery)
	if err != nil {
		log.Println(err)
		return
	}
	collection := d.deliveryDB.GetDatabase("ts").GetCollection("delivery")
	err = collection.InsertOne(delivery)
	if err != nil {
		log.Println(err)
	}
}

func (d *deliveryService) Entry() {
	d.deliveryQueue.Recv(d.ProcessDelivery)
}
