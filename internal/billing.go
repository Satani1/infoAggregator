package internal

import (
	"finalWork/internal/models"
	"io/ioutil"
	"math"
	"strconv"
)

func (app *Application) Billing() (BillingData *models.BillingData, err error) {
	nums, err := ioutil.ReadFile("./data/billing.data")
	if err != nil {
		return BillingData, err
	}
	var Bits []int

	for i := 5; i >= 0; i-- {
		byteNumber, _ := strconv.Atoi(string(nums[i]))
		Bits = append(Bits, byteNumber)
	}
	var sum uint8
	for index, num := range Bits {
		if num == 1 {
			sum += uint8(math.Pow(float64(2), float64(index)))
		}
	}

	BillingData = &models.BillingData{

		CreateCustomer: BillCheck(nums[0]),
		Purchase:       BillCheck(nums[1]),
		Payout:         BillCheck(nums[2]),
		Recurring:      BillCheck(nums[3]),
		FraudControl:   BillCheck(nums[4]),
		CheckoutPage:   BillCheck(nums[5]),
	}

	return BillingData, nil
}
func BillCheck(bit uint8) bool {
	if int(bit) == 49 {
		return true
	} else {
		return false
	}
}
