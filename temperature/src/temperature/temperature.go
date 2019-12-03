package temperature

func ConvertFtoC(f float64) float64{
    var c float64
    c = (f-32)/1.8
    return c
}

func ConvertCtoF(c float64) float64{
    var f float64
    f = c*1.8+32
    return f
}
