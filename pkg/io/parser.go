package io

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/rodoufu/simple-orderbook/pkg/entity"
)

func ReadTransactions(ctx context.Context, fileName string) (<-chan Transaction, error) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "problem opening: %v", fileName)
	}

	resp := make(chan Transaction)
	go func() {
		defer csvFile.Close()
		defer close(resp)

		done := ctx.Done()
		csvReader := csv.NewReader(csvFile)
		csvReader.Comment = '#'

		for {
			select {
			case <-done:
				return
			default:
				var record []string
				record, err = csvReader.Read()
				if err == io.EOF {
					return
				}

				for i := 0; i < len(record); i++ {
					record[i] = strings.TrimSpace(record[i])
				}

				switch record[0] {
				case "N":
					if len(record) != 7 {
						resp <- ErrorTransaction{
							Err: fmt.Errorf("invalid create order line: %v", record),
						}
						return
					}

					var userID, price, amount, orderID int64
					userID, err = strconv.ParseInt(record[1], 10, 64)
					if err != nil {
						resp <- ErrorTransaction{
							Err: errors.Wrapf(err, "problem parsing user ID in create order"),
						}
						return
					}
					price, err = strconv.ParseInt(record[3], 10, 64)
					if err != nil {
						resp <- ErrorTransaction{
							Err: errors.Wrapf(err, "problem parsing price in create order"),
						}
						return
					}
					amount, err = strconv.ParseInt(record[4], 10, 64)
					if err != nil {
						resp <- ErrorTransaction{
							Err: errors.Wrapf(err, "problem parsing amount in create order"),
						}
						return
					}
					orderID, err = strconv.ParseInt(record[6], 10, 64)
					if err != nil {
						resp <- ErrorTransaction{
							Err: errors.Wrapf(err, "problem parsing amount in create order"),
						}
						return
					}

					side := entity.Buy
					if record[5] == "S" {
						side = entity.Sell
					}
					resp <- NewOrderTransaction{
						Symbol: record[2],
						Order: entity.Order{
							Amount:    uint64(amount),
							Price:     uint64(price),
							ID:        entity.OrderID(orderID),
							Side:      side,
							User:      entity.UserID(userID),
							Timestamp: time.Now(),
						},
					}
				case "C":
					if len(record) != 3 {
						resp <- ErrorTransaction{
							Err: fmt.Errorf("invalid cancel order line: %v", record),
						}
						return
					}

					var userID, orderID int64
					userID, err = strconv.ParseInt(record[1], 10, 64)
					if err != nil {
						resp <- ErrorTransaction{
							Err: errors.Wrapf(err, "problem parsing user ID in cancel order"),
						}
						return
					}
					orderID, err = strconv.ParseInt(record[2], 10, 64)
					if err != nil {
						resp <- ErrorTransaction{
							Err: errors.Wrapf(err, "problem parsing order ID in cancel order"),
						}
						return
					}
					resp <- CancelOrderTransaction{
						User:    entity.UserID(userID),
						OrderID: entity.OrderID(orderID),
					}
				case "F":
					resp <- FlushAllOrdersTransaction{}
				default:
					resp <- ErrorTransaction{
						Err: fmt.Errorf("invalid line: %v", record),
					}
				}
			}
		}
	}()
	return resp, nil
}
