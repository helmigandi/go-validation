package go_validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestValidation(t *testing.T) {
	validate := validator.New()
	if validate == nil {
		t.Error("validation is nil")
	}
}

func TestValidationVariable(t *testing.T) {
	validate := validator.New()
	username := "andi"

	err := validate.Var(username, "required,alphanum")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationCompareTwoVariable(t *testing.T) {
	validate := validator.New()
	password := "password"
	confirmPassword := "password"

	err := validate.VarWithValue(password, confirmPassword, "eqfield")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationParameter(t *testing.T) {
	validate := validator.New()
	phoneNumber := "0815900141"

	err := validate.Var(phoneNumber, "required,numeric,min=5,max=10")
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationStruct(t *testing.T) {
	type LoginRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=5,max=64"`
	}

	validate := validator.New()
	err := validate.Struct(LoginRequest{Username: "andi@mail.com", Password: "password"})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationErrors(t *testing.T) {
	type LoginRequest struct {
		Username string `json:"username" validate:"required,email"`
		Password string `json:"password" validate:"required,min=5,max=64"`
		FullName string `json:"full_name" validate:"required,alphanum"`
	}

	validate := validator.New()
	err := validate.Struct(LoginRequest{Username: "andi", Password: "asd", FullName: "andi soraya@"})
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			fmt.Println("Field: ", e.Field(), ", Tage: ", e.Tag(), ", Value: ", e.Value(), ", Error: ", e.Error())
		}
	}
}

func TestValidationStructCross(t *testing.T) {
	type RegisterRequest struct {
		Username        string `json:"username" validate:"required,email"`
		Password        string `json:"password" validate:"required,min=5,max=64"`
		ConfirmPassword string `json:"confirm_password" validate:"required,min=5,max=64,eqfield=Password"`
	}

	validate := validator.New()
	err := validate.Struct(RegisterRequest{Username: "andi@mail.com", Password: "password", ConfirmPassword: "password"})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationStructNested(t *testing.T) {
	type AddressRequest struct {
		Street string `json:"street" validate:"required,alphanum"`
		City   string `json:"city" validate:"required,alphanum"`
	}

	type RegisterRequest struct {
		Username        string         `json:"username" validate:"required,email"`
		Password        string         `json:"password" validate:"required,min=5,max=64"`
		ConfirmPassword string         `json:"confirm_password" validate:"required,min=5,max=64,eqfield=Password"`
		Address         AddressRequest `json:"address"`
	}

	validate := validator.New()
	err := validate.Struct(RegisterRequest{
		Username: "andi@mail.com", Password: "password", ConfirmPassword: "password",
		Address: AddressRequest{Street: "Test", City: "TestCity"},
	})

	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationStructNestedSlice(t *testing.T) {
	type AddressRequest struct {
		Street string `json:"street" validate:"required,alphanum"`
		City   string `json:"city" validate:"required,alphanum"`
	}

	type RegisterRequest struct {
		Username        string           `json:"username" validate:"required,email"`
		Password        string           `json:"password" validate:"required,min=5,max=64"`
		ConfirmPassword string           `json:"confirm_password" validate:"required,min=5,max=64,eqfield=Password"`
		Addresses       []AddressRequest `json:"addresses" validate:"required,dive"`
	}

	validate := validator.New()
	err := validate.Struct(RegisterRequest{
		Username: "andi@mail.com", Password: "password", ConfirmPassword: "password",
		Addresses: []AddressRequest{{Street: "Test", City: "TestCity"}},
	})

	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationBasicCollection(t *testing.T) {
	type LoginRequest struct {
		Username string   `json:"username" validate:"required,email"`
		Password string   `json:"password" validate:"required,min=5,max=64"`
		Hobbies  []string `json:"hobbies" validate:"required,dive,alphanum,min=1"`
	}

	validate := validator.New()
	err := validate.Struct(LoginRequest{Username: "andi@mail.com", Password: "password", Hobbies: []string{"Coding"}})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationMap(t *testing.T) {
	type School struct {
		Name string `validate:"required,min=8"`
	}

	type LoginRequest struct {
		Username string            `json:"username" validate:"required,email"`
		Password string            `json:"password" validate:"required,min=5,max=64"`
		Schools  map[string]School `validate:"dive,keys,required,min=2,endkeys,required"`
	}

	validate := validator.New()
	err := validate.Struct(LoginRequest{
		Username: "andi@mail.com", Password: "password",
		Schools: map[string]School{"SMA": {Name: "SMA 1 Bekasi"}}})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestValidationBasicMap(t *testing.T) {

	type LoginRequest struct {
		Username string         `json:"username" validate:"required,email"`
		Password string         `json:"password" validate:"required,min=5,max=64"`
		Scores   map[string]int `json:"scores" validate:"dive,keys,required,min=2,endkeys,required,gt=0"`
	}

	validate := validator.New()
	err := validate.Struct(LoginRequest{
		Username: "andi@mail.com", Password: "password",
		Scores: map[string]int{"Math": 10},
	})
	if err != nil {
		t.Error(err.Error())
	}
}

func TestAlias(t *testing.T) {
	validate := validator.New()
	validate.RegisterAlias("varchar", "required,max=255")

	type Seller struct {
		Id     string `validate:"varchar,min=5"`
		Name   string `validate:"varchar"`
		Owner  string `validate:"varchar"`
		Slogan string `validate:"varchar"`
	}

	seller := Seller{
		Id:     "123",
		Name:   "",
		Owner:  "",
		Slogan: "",
	}

	err := validate.Struct(seller)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func MustValidUsername(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if ok {
		if value != strings.ToUpper(value) {
			return false
		}
		if len(value) < 5 {
			return false
		}
	}
	return true
}

func TestCustomValidationFunction(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("username", MustValidUsername)

	type LoginRequest struct {
		Username string `validate:"required,username"`
		Password string `validate:"required"`
	}

	request := LoginRequest{
		Username: "ABDUL",
		Password: "",
	}

	err := validate.Struct(request)
	if err != nil {
		fmt.Println(err.Error())
	}
}

var regexNumber = regexp.MustCompile("^[0-9]+$")

func MustValidPin(field validator.FieldLevel) bool {
	length, err := strconv.Atoi(field.Param())
	if err != nil {
		panic(err)
	}

	value := field.Field().String()
	if !regexNumber.MatchString(value) {
		return false
	}

	return len(value) == length
}

func TestCustomValidationParameter(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("pin", MustValidPin)

	type Login struct {
		Phone string `validate:"required,number"`
		Pin   string `validate:"required,pin=6"`
	}

	request := Login{
		Phone: "0904190424",
		Pin:   "123123",
	}

	err := validate.Struct(request)
	if err != nil {
		fmt.Println(err)
	}
}

func TestOrRule(t *testing.T) {
	type Login struct {
		Username string `validate:"required,email|numeric"`
		Password string `validate:"required"`
	}

	request := Login{
		Username: "12345",
		Password: "ekoo",
	}

	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func MustEqualsIgnoreCase(field validator.FieldLevel) bool {
	value, _, _, ok := field.GetStructFieldOK2()
	if !ok {
		panic("field not ok")
	}

	firstValue := strings.ToUpper(field.Field().String())
	secondValue := strings.ToUpper(value.String())

	return firstValue == secondValue
}

func TestCrossFieldValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("field_equals_ignore_case", MustEqualsIgnoreCase)

	type User struct {
		// username harus sama dengan email atau phone.
		Username string `validate:"required,field_equals_ignore_case=Email|field_equals_ignore_case=Phone"`
		Email    string `validate:"required,email"`
		Phone    string `validate:"required,numeric"`
		Name     string `validate:"required"`
	}

	user := User{
		Username: "eko@example.com",
		Email:    "eko@example.com",
		Phone:    "089999999999",
		Name:     "Eko",
	}

	err := validate.Struct(user)
	if err != nil {
		fmt.Println(err)
	}
}

type RegisterRequest struct {
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Phone    string `validate:"required,numeric"`
	Password string `validate:"required"`
}

func MustValidRegisterSuccess(level validator.StructLevel) {
	registerRequest := level.Current().Interface().(RegisterRequest)

	if registerRequest.Username == registerRequest.Email || registerRequest.Username == registerRequest.Phone {
		// sukses
	} else {
		// gagal
		level.ReportError(registerRequest.Username, "Username", "Username", "username", "")
	}
}

func TestStructLevelValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(MustValidRegisterSuccess, RegisterRequest{})

	request := RegisterRequest{
		Username: "089923942934",
		Email:    "eko@example.com",
		Phone:    "089923942934",
		Password: "rahasia",
	}

	err := validate.Struct(request)
	if err != nil {
		fmt.Println(err.Error())
	}
}
