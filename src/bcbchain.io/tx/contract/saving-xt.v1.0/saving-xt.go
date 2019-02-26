package saving_xt_v1_0

import (
	atm "bcbchain.io/algorithm"
	"bcbchain.io/prototype"
	"bcbchain.io/tx"
	tx2 "bcbchain.io/tx/tx"
	"bcbchain.io/utils"
)

func SxtInit(nonce, gasLimit, note, smcAddress, privateKey string) string {
	items := make([]string, 0)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtInit)),
		items,
		privateKey,
	)
}

func SxtActivate(nonce, gasLimit, note, smcAddress, privateKey string) string {
	items := make([]string, 0)
	return tx.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtActivate)),
		items,
		privateKey,
	)
}

func SxtSetOwner(nonce, gasLimit, note, smcAddress, newOwner, privateKey string) string {
	items := make([]string, 0)
	items = append(items, newOwner)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtActivate)),
		items,
		privateKey,
	)
}

func SxtRegisterNameXid(nonce, gasLimit, note, smcAddress, name, affCode, bcb, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, affCode)
	items = append(items, bcb)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtRegisterNameXid)),
		items,
		privateKey,
	)
}

func SxtRegisterNameXaddr(nonce, gasLimit, note, smcAddress, name, affCode, bcb, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, affCode)
	items = append(items, bcb)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtRegisterNameXaddr)),
		items,
		privateKey,
	)
}

func SxtRegisterNameXname(nonce, gasLimit, note, smcAddress, name, affCode, bcb, privateKey string) string {
	items := make([]string, 0)
	items = append(items, name)
	items = append(items, affCode)
	items = append(items, bcb)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtRegisterNameXname)),
		items,
		privateKey,
	)
}

func SxtBuy(nonce, gasLimit, note, smcAddress, team, bcb, privateKey string) string {
	items := make([]string, 0)
	items = append(items, team)
	items = append(items, bcb)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtBuy)),
		items,
		privateKey,
	)
}

func SxtBuyXid(nonce, gasLimit, note, smcAddress, team, bcb, affCode, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, bcb)

	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtBuyXid)),
		items,
		privateKey,
	)
}

func SxtBuyXAddress(nonce, gasLimit, note, smcAddress, team, bcb, affAddr, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affAddr)
	items = append(items, team)
	items = append(items, bcb)

	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtBuyXaddr)),
		items,
		privateKey,
	)
}

func SxtBuyXName(nonce, gasLimit, note, smcAddress, team, bcb, affName, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affName)
	items = append(items, team)
	items = append(items, bcb)

	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtBuyXname)),
		items,
		privateKey,
	)
}

func SxtReloadXid(nonce, gasLimit, note, smcAddress, team, bcb, affCode, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, bcb)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtReloadXid)),
		items,
		privateKey,
	)
}

func SxtReloadXaddr(nonce, gasLimit, note, smcAddress, team, bcb, affCode, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, bcb)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtReloadXaddr)),
		items,
		privateKey,
	)
}

func SxtReloadXname(nonce, gasLimit, note, smcAddress, team, bcb, affCode, privateKey string) string {
	items := make([]string, 0)
	items = append(items, affCode)
	items = append(items, team)
	items = append(items, bcb)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtReloadXname)),
		items,
		privateKey,
	)
}

func SxtWithDraw(nonce, gasLimit, note, smcAddress, privateKey string) string {
	items := make([]string, 0)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtWithDraw)),
		items,
		privateKey,
	)
}

func SxtSetSettings(nonce, gasLimit, note, smcAddress, settings, privateKey string) string {
	items := make([]string, 0)
	items = append(items, settings)
	return tx2.PackAndSignTx(
		nonce,
		gasLimit,
		note,
		smcAddress,
		utils.BytesToHex(atm.CalcMethodId(prototype.SxtSetSetting)),
		items,
		privateKey,
	)
}
