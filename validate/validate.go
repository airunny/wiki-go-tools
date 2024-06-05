package validate

import "regexp"

var (
	urlReg      = regexp.MustCompile(`^(https|http)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)
	domainReg   = regexp.MustCompile(`^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+(?:[a-zA-Z]{2,})$`)
	mobileReg   = regexp.MustCompile(`^1[3456789]\d{9}$`)
	emailReg    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	wechatReg   = regexp.MustCompile(`^[a-zA-Z0-9_-]{6,20}$`)
	qqReg       = regexp.MustCompile(`^[1-9]\d{4,10}$`)
	landlineReg = regexp.MustCompile(`^\d{7,8}$`)
	numericReg  = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	floatReg    = regexp.MustCompile(`^[0-9]+(\.[0-9]{0,5})?$`)
)

func ValidURL(in string) bool {
	return urlReg.MatchString(in)
}

func ValidDomain(in string) bool {
	return domainReg.MatchString(in)
}

func ValidMobile(in string) bool {
	return mobileReg.MatchString(in)
}

func ValidEmail(in string) bool {
	return emailReg.MatchString(in)
}

func ValidWechat(in string) bool {
	return wechatReg.MatchString(in)
}

func ValidQQ(in string) bool {
	return qqReg.MatchString(in)
}

func ValidLandline(in string) bool {
	return landlineReg.MatchString(in)
}

func ValidNumeric(in string) bool {
	return numericReg.MatchString(in)
}

func ValidFloat(in string) bool {
	return floatReg.MatchString(in)
}
