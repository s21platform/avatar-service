# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/avatar.proto](#api_avatar-proto)
    - [Avatar](#-Avatar)
    - [DeleteSocietyAvatarIn](#-DeleteSocietyAvatarIn)
    - [DeleteUserAvatarIn](#-DeleteUserAvatarIn)
    - [GetAllSocietyAvatarsIn](#-GetAllSocietyAvatarsIn)
    - [GetAllSocietyAvatarsOut](#-GetAllSocietyAvatarsOut)
    - [GetAllUserAvatarsOut](#-GetAllUserAvatarsOut)
    - [NewAvatarRegister](#-NewAvatarRegister)
    - [SetSocietyAvatarIn](#-SetSocietyAvatarIn)
    - [SetSocietyAvatarOut](#-SetSocietyAvatarOut)
    - [SetUserAvatarIn](#-SetUserAvatarIn)
    - [SetUserAvatarOut](#-SetUserAvatarOut)
  
    - [AvatarService](#-AvatarService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_avatar-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/avatar.proto



<a name="-Avatar"></a>

### Avatar



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int32](#int32) |  |  |
| link | [string](#string) |  |  |






<a name="-DeleteSocietyAvatarIn"></a>

### DeleteSocietyAvatarIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| avatar_id | [int32](#int32) |  |  |






<a name="-DeleteUserAvatarIn"></a>

### DeleteUserAvatarIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| avatar_id | [int32](#int32) |  |  |






<a name="-GetAllSocietyAvatarsIn"></a>

### GetAllSocietyAvatarsIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |






<a name="-GetAllSocietyAvatarsOut"></a>

### GetAllSocietyAvatarsOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| avatar_list | [Avatar](#Avatar) | repeated |  |






<a name="-GetAllUserAvatarsOut"></a>

### GetAllUserAvatarsOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| avatar_list | [Avatar](#Avatar) | repeated |  |






<a name="-NewAvatarRegister"></a>

### NewAvatarRegister
kafka contracts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| link | [string](#string) |  |  |






<a name="-SetSocietyAvatarIn"></a>

### SetSocietyAvatarIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| filename | [string](#string) |  |  |
| batch | [bytes](#bytes) |  |  |






<a name="-SetSocietyAvatarOut"></a>

### SetSocietyAvatarOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| link | [string](#string) |  |  |






<a name="-SetUserAvatarIn"></a>

### SetUserAvatarIn



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| filename | [string](#string) |  |  |
| batch | [bytes](#bytes) |  |  |






<a name="-SetUserAvatarOut"></a>

### SetUserAvatarOut



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| link | [string](#string) |  |  |





 

 

 


<a name="-AvatarService"></a>

### AvatarService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SetUserAvatar | [.SetUserAvatarIn](#SetUserAvatarIn) stream | [.SetUserAvatarOut](#SetUserAvatarOut) |  |
| GetAllUserAvatars | [.google.protobuf.Empty](#google-protobuf-Empty) | [.GetAllUserAvatarsOut](#GetAllUserAvatarsOut) |  |
| DeleteUserAvatar | [.DeleteUserAvatarIn](#DeleteUserAvatarIn) | [.Avatar](#Avatar) |  |
| SetSocietyAvatar | [.SetSocietyAvatarIn](#SetSocietyAvatarIn) stream | [.SetSocietyAvatarOut](#SetSocietyAvatarOut) |  |
| GetAllSocietyAvatars | [.GetAllSocietyAvatarsIn](#GetAllSocietyAvatarsIn) | [.GetAllSocietyAvatarsOut](#GetAllSocietyAvatarsOut) |  |
| DeleteSocietyAvatar | [.DeleteSocietyAvatarIn](#DeleteSocietyAvatarIn) | [.Avatar](#Avatar) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

