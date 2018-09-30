package amount

import (
	"testing"
)

func TestFromNumString(t *testing.T) {
	// t.SkipNow()
	test(t, "100.001", "10000100000", "100.00100000", "100.001")
	test(t, ".001", "00100000", "0.00100000", "0.001")
	test(t, "0.001", "00100000", "0.00100000", "0.001")
	test(t, "100", "10000000000", "100.0", "100.0")
	test(t, "0", "00000000", "0.0", "0.0")
	test(t, "-0", "00000000", "0.0", "0.0")
	test(t, ".00000001", "00000001", "0.00000001", "0.00000001")
	test(t, "-100.001", "-10000100000", "-100.00100000", "-100.001")
	test(t, "-100", "-10000000000", "-100.0", "-100.0")
	test(t, "-.001", "-00100000", "-0.00100000", "-0.001")
	test(t, "-0.001", "-00100000", "-0.00100000", "-0.001")
	test(t, "-0010.001", "-1000100000", "-10.00100000", "-10.001")
	test(t, "10.001", "1000100000", "10.00100000", "10.001")
	test(t, "-10.001", "-1000100000", "-10.00100000", "-10.001")
	test(t, "", "00000000", "0.0", "0.0")
	test(t, "123456789.98765432", "12345678998765432", "123456789.98765432", "123456789.98765432")
	test(t, "-123456789.98765432", "-12345678998765432", "-123456789.98765432", "-123456789.98765432")
	test(t, "023456789.98765430", "2345678998765430", "23456789.98765430", "23456789.9876543")
	test(t, "-023456789.98765430", "-2345678998765430", "-23456789.98765430", "-23456789.9876543")
}

func test(t *testing.T, value string, result string, floatStr string, floatShort string) {
	amount, err := FromNumString(value)
	if err != nil {
		t.Errorf("FromNumString Error:%s", err)
		return
	}
	if amount != Amount(result) {
		t.Errorf("Amount:[%s] %s != result(%s)", value, amount, result)
	}
	floatStrR := amount.FloatString()
	if floatStrR != floatStr {
		t.Errorf("FloatString: %s != %s", floatStrR, floatStr)
	}
	floatShortR := amount.FloatStringShort()
	if floatShortR != floatShort {
		t.Errorf("FoatStringShort: %s != %s", floatShortR, floatShort)
	}
}

func TestMinus(t *testing.T) {
	a, _ := FromNumString("-1.1")
	aMinus := a.Minus()
	if aMinus != "-110000000" {
		t.Errorf("Minus:[-1.1] %s %s", a, aMinus)
	}
	a, _ = FromNumString("1.1")
	aMinus = a.Minus()
	if aMinus != "-110000000" {
		t.Errorf("Minus:[1.1] %s %s", a, aMinus)
	}
	a, _ = FromNumString("0")
	aMinus = a.Minus()
	if aMinus != "00000000" {
		t.Errorf("Minus:[1.1] %s %s", a, aMinus)
	}
}

func TestAdd(t *testing.T) {
	testAdd(t, "1", "2", "3")
	testAdd(t, "2", "2", "4")
	testAdd(t, "1.1", "1.1", "2.2")
	testAdd(t, "1.11111111", "1.11111111", "2.22222222")
	testAdd(t, "-1", "1", "0")
	testAdd(t, "-1.1", "1.1", "0")
	testAdd(t, "0", "0", "0")
}

func testAdd(t *testing.T, x, y, z string) {
	amountX, _ := FromNumString(x)
	amountY, _ := FromNumString(y)
	amountZ, _ := FromNumString(z)
	addXY := amountX.Add(amountY)
	if addXY != amountZ {
		t.Errorf("Add:[%s + %s = %s] != [%s + %s = %s(%s)]", x, y, z, amountX, amountY, addXY, amountZ)
	}
}

func TestMul(t *testing.T) {
	testMul(t, "1", "1", "1")
	testMul(t, "2", "3", "6")
	testMul(t, "2.1", "3.1", "6.51")
	testMul(t, "0.1", "0.1", "0.01")
	testMul(t, "0.0000001", "0.000000001", "0")
	testMul(t, "0.11111111", "0.11111111", "0.01234567")
	testMul(t, "1.11111111", "1.11111111", "1.23456789")
	testMul(t, "11111111", "11111111", "123456787654321")

}

func testMul(t *testing.T, x, y, z string) {
	amountX, _ := FromNumString(x)
	amountY, _ := FromNumString(y)
	amountZ, _ := FromNumString(z)
	addXY := amountX.Mul(amountY)
	if addXY != amountZ {
		t.Errorf("Add:[%s * %s = %s] != [%s * %s = %s(%s)]", x, y, z, amountX, amountY, addXY, amountZ)
	}
}

func TestDiv(t *testing.T) {
	testDiv(t, "1", "1", "1")
	testDiv(t, "0.1", "0.1", "1")
	testDiv(t, "0.11", "0.1", "1.1")
	testDiv(t, "0.009", "100", "0.00009")
	testDiv(t, "0.00000009", "100", "0")
	testDiv(t, "1", "3", "0.33333333")
}

func testDiv(t *testing.T, x, y, z string) {
	amountX, _ := FromNumString(x)
	amountY, _ := FromNumString(y)
	amountZ, _ := FromNumString(z)
	addXY := amountX.Div(amountY)
	if addXY != amountZ {
		t.Errorf("Add:[%s / %s = %s] != [%s / %s = %s(%s)]", x, y, z, string(amountX), string(amountY), string(addXY), string(amountZ))
	}
}
