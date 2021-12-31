package mth

//求最大公约数
func GreaterCommonDivisor(a,b int)int{
	if a % b == 0{
		return b
	}
	return GreaterCommonDivisor(b ,a %b)
}
