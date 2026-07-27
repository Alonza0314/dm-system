# SettingApi

All URIs are relative to *http://127.0.0.1:5000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**settingAccount**](#settingaccount) | **POST** /api/setting/account | Update account settings|

# **settingAccount**
> settingAccount(settingAccountRequest)


### Example

```typescript
import {
    SettingApi,
    Configuration,
    SettingAccountRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new SettingApi(configuration);

let settingAccountRequest: SettingAccountRequest; //

const { status, data } = await apiInstance.settingAccount(
    settingAccountRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **settingAccountRequest** | **SettingAccountRequest**|  | |


### Return type

void (empty response body)

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
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

