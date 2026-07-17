# QrcodeApi

All URIs are relative to *http://127.0.0.1:5000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**borrowDevice**](#borrowdevice) | **POST** /api/qrcode/{category}/{name} | Borrow device|
|[**returnDevice**](#returndevice) | **DELETE** /api/qrcode/{category}/{name} | Return device|

# **borrowDevice**
> borrowDevice(deviceUserRequest)


### Example

```typescript
import {
    QrcodeApi,
    Configuration,
    DeviceUserRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new QrcodeApi(configuration);

let category: string; //Category name (default to undefined)
let name: string; //Device name (default to undefined)
let deviceUserRequest: DeviceUserRequest; //

const { status, data } = await apiInstance.borrowDevice(
    category,
    name,
    deviceUserRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceUserRequest** | **DeviceUserRequest**|  | |
| **category** | [**string**] | Category name | defaults to undefined|
| **name** | [**string**] | Device name | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **returnDevice**
> returnDevice(deviceUserRequest)


### Example

```typescript
import {
    QrcodeApi,
    Configuration,
    DeviceUserRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new QrcodeApi(configuration);

let category: string; //Category name (default to undefined)
let name: string; //Device name (default to undefined)
let deviceUserRequest: DeviceUserRequest; //

const { status, data } = await apiInstance.returnDevice(
    category,
    name,
    deviceUserRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceUserRequest** | **DeviceUserRequest**|  | |
| **category** | [**string**] | Category name | defaults to undefined|
| **name** | [**string**] | Device name | defaults to undefined|


### Return type

void (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | Not Found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

