# HistoryApi

All URIs are relative to *http://127.0.0.1:5000*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getHistory**](#gethistory) | **GET** /api/history/{category}/{name} | Get history|

# **getHistory**
> HistoryListResponse getHistory()


### Example

```typescript
import {
    HistoryApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new HistoryApi(configuration);

let category: string; //Category name (default to undefined)
let name: string; //Device name (default to undefined)

const { status, data } = await apiInstance.getHistory(
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

**HistoryListResponse**

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

