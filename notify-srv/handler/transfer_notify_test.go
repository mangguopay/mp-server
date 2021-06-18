package handler

import "testing"

func TestTransferNotifyHandler_TransferSuccessNotify(t1 *testing.T) {
	TransferNotifyH.TransferSuccessNotify("2020102916403620769039")
}
