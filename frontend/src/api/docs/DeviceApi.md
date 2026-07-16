# DeviceApi

All URIs are relative to *http://127.0.0.1:5000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createDevice**](#createdevice) | **POST** /api/device | Create device|
|[**deleteDevice**](#deletedevice) | **DELETE** /api/device/{category}/{name} | Delete device|
|[**getDevice**](#getdevice) | **GET** /api/device/{category}/{name} | Get device|
|[**listDevices**](#listdevices) | **GET** /api/device/{category} | List devices|

# **createDevice**
> MessageResponse createDevice(deviceCreateRequest)


### Example

```typescript
import {
    DeviceApi,
    Configuration,
    DeviceCreateRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new DeviceApi(configuration);

let deviceCreateRequest: DeviceCreateRequest; //

const { status, data } = await apiInstance.createDevice(
    deviceCreateRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **deviceCreateRequest** | **DeviceCreateRequest**|  | |


### Return type

**MessageResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteDevice**
> MessageResponse deleteDevice()


### Example

```typescript
import {
    DeviceApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeviceApi(configuration);

let category: string; //Category name (default to undefined)
let name: string; //Device name (default to undefined)

const { status, data } = await apiInstance.deleteDevice(
    category,
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **category** | [**string**] | Category name | defaults to undefined|
| **name** | [**string**] | Device name | defaults to undefined|


### Return type

**MessageResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**409** | Conflict |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getDevice**
> Device getDevice()


### Example

```typescript
import {
    DeviceApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeviceApi(configuration);

let category: string; //Category name (default to undefined)
let name: string; //Device name (default to undefined)

const { status, data } = await apiInstance.getDevice(
    category,
    name
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **category** | [**string**] | Category name | defaults to undefined|
| **name** | [**string**] | Device name | defaults to undefined|


### Return type

**Device**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listDevices**
> DeviceListResponse listDevices()


### Example

```typescript
import {
    DeviceApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DeviceApi(configuration);

let category: string; //Category name (default to undefined)

const { status, data } = await apiInstance.listDevices(
    category
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **category** | [**string**] | Category name | defaults to undefined|


### Return type

**DeviceListResponse**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

