package resp

import (
	"fmt"
	"strconv"
)

func GetByteTillClrf(input []byte) ([]byte, []byte, bool) {
	if len(input) == 0 {
		return []byte{}, []byte{}, false
	}
	index := 0

	for index+1 < len(input) && input[index] != '\r' && input[index+1] != '\n' {
		index++
	}

	if index+1 < len(input) && input[index] == '\r' && input[index+1] == '\n' {
		return input[:index], input[index+2:], true
	}
	return input, []byte{}, false

}

type Resp interface {
	AsBytes() []byte
}

type respSimpleString struct {
	inner string
}

func (rs *respSimpleString) AsBytes() []byte {
	result := make([]byte, 0, 1+len(rs.inner))
	result = append(result, []byte("+")...)
	result = append(result, []byte(rs.inner)...)
	result = append(result, []byte("\r\n")...)
	return result
}

type respErrorString struct {
	inner string
}

func (rs *respErrorString) AsBytes() []byte {
	result := make([]byte, 0, 1+len(rs.inner))
	result = append(result, []byte("-")...)
	result = append(result, []byte(rs.inner)...)
	result = append(result, []byte("\r\n")...)
	return result
}

type respInteger struct {
	inner int
}

func (rs *respInteger) AsBytes() []byte {
	integertostring := fmt.Sprintf("%d",rs.inner)
	result := make([]byte, 0, 1+len(integertostring))
	result = append(result, []byte(":")...)
	result = append(result, []byte(integertostring)...)
	result = append(result, []byte("\r\n")...)
	return result
}

type respBulkString struct {
	inner string
}	
func (rs *respBulkString) AsBytes()[]byte{
	length:=fmt.Sprintf("%d",len(rs.inner));
	bulkStringToByte:=[]byte(rs.inner);
	result:=make([]byte,1+len(length)+2+len(bulkStringToByte)+2);
	result=append(result, []byte("$")...);
	result=append(result, []byte(length)...);
	result=append(result, []byte("\r\n")...);
	result=append(result, []byte(bulkStringToByte)...);
	result=append(result, []byte("\r\n")...);
	return result;

}

type respArray struct{
	inner []Resp
}
func (rs *respArray) AsBytes()[]byte{
	length:=fmt.Sprintf("%d",len(rs.inner));
	
	result:=make([]byte,1+len(length)+2+len(rs.inner)+2);
	result=append(result, []byte("*")...);
	result=append(result, []byte(length)...);
	result=append(result, []byte("\r\n")...);
	for _,val:= range rs.inner {
		result=append(result,val.AsBytes()...);
	}
	result = append(result, []byte("\r\n")...);
	return result;
}

type respNilArray struct{
	
}
func (rs *respNilArray)AsBytes()[]byte{
	return []byte("*-1\r\n");
}

type respNilString struct{}

func (rs* respNilString)AsBytes()[]byte{
	return []byte("$-1\r\n")
}

func RespFromBytes(input []byte) ([]byte, Resp) {
	if len(input) == 0 {
		return []byte{}, nil
	}
	if input[0] == '+' {
		str, leftovers, valid := GetByteTillClrf(input[1:])
		if !valid {
			return []byte{}, nil
		}
		rs := respSimpleString{
			inner: string(str),
		}
		return leftovers, &rs
	}else if input[0]=='-'{
		str, leftovers, valid := GetByteTillClrf(input[1:])
		if !valid {
			return []byte{}, nil
		}
		rs := respErrorString{
			inner: string(str),
		}
		return leftovers, &rs
	}else if input[0]==':' {
		str,leftover,valid:=GetByteTillClrf(input[1:]);
		if(!valid){
			return []byte{},nil;
		}
		val,err:=strconv.ParseInt(string(str),10,32);
		if(err!=nil){return []byte{},nil }
		rs:=respInteger{inner: int(val)}
		return leftover,&rs

	}else if input[0]=='$'{
		//bulk string
		strLength,leftover,valid:=GetByteTillClrf(input[1:]);
		if(!valid){
			return []byte{},nil;
		}
		len64,err:=strconv.ParseInt(string(strLength),10,32);
		len32:=int(len64)
		if(err!=nil){
			return []byte{},nil;
		}

		if len32<0{
			return leftover,&respNilString{}
		}

		if len32+1<len(leftover)&&leftover[len32]=='\r'&&leftover[len32+1]=='\n'{
			rs:=respBulkString{
				inner: string(leftover),
			}
			return leftover,&rs;
		}
	}else if input[0]=='*'{
		//array
		arrLength,inputArray,valid:=GetByteTillClrf(input[1:]);
		if(!valid){
			return []byte{},nil;
		}
		len64,err:=strconv.ParseInt(string(arrLength),10,32);
		if(err!=nil){
			return []byte{},nil;
		}
		len32:=int(len64);
		if(len32<0){return inputArray,&respNilArray{}} ;
		
		res:=make([]Resp,0,len32)

		for i:=0;i<len32&&len(inputArray)!=0;i++{
			leftovers,currentArray:=RespFromBytes(inputArray);
			if(currentArray==nil){
				return []byte{},nil;
			}
			inputArray=leftovers;
			res = append(res, currentArray);
		}

		if(len(res)!=len32){
			return []byte{},nil;
		}
		rs:=respArray{
			inner: res,
		}
		return inputArray,&rs;
		}else{
			return []byte{},nil;
		}

		return []byte{}, nil
	}


