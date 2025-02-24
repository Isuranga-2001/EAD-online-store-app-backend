package com.ead.payment_service.helper;

import com.ead.payment_service.exception.RequestFailedException;
import org.apache.hc.client5.http.classic.methods.*;
import org.apache.hc.client5.http.impl.classic.CloseableHttpClient;
import org.apache.hc.client5.http.impl.classic.CloseableHttpResponse;
import org.apache.hc.client5.http.impl.classic.HttpClients;
import org.apache.hc.core5.http.io.entity.StringEntity;

public class RequestHelper {
    
    public static CloseableHttpResponse SendGetRequest(String url) {
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpGet httpGet = new HttpGet(url);
            return httpClient.execute(httpGet);
        } catch (Exception e) {
            throw new RequestFailedException("GET request failed for URL: " + url, e);
        }
    }

    public static int SendPostRequest(String url, String body) {
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpPost httpPost = new HttpPost(url);
            httpPost.setEntity(new StringEntity(body));
            try (var response = httpClient.execute(httpPost)) {
                return response.getCode();
            }
        } catch (Exception e) {
            throw new RequestFailedException("POST request failed for URL: " + url, e);
        }
    }

    public static int SendPutRequest(String url, String body) {
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpPut httpPut = new HttpPut(url);
            httpPut.setEntity(new StringEntity(body));
            try (var response = httpClient.execute(httpPut)) {
                return response.getCode();
            }
        } catch (Exception e) {
            throw new RequestFailedException("PUT request failed for URL: " + url, e);
        }
    }

    public static int SendDeleteRequest(String url) {
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpDelete httpDelete = new HttpDelete(url);
            try (var response = httpClient.execute(httpDelete)) {
                return response.getCode();
            }
        } catch (Exception e) {
            throw new RequestFailedException("DELETE request failed for URL: " + url, e);
        }
    }

    public static int SendPatchRequest(String url, String body) {
        try (CloseableHttpClient httpClient = HttpClients.createDefault()) {
            HttpPatch httpPatch = new HttpPatch(url);
            httpPatch.setEntity(new StringEntity(body));
            try (var response = httpClient.execute(httpPatch)) {
                return response.getCode();
            }
        } catch (Exception e) {
            throw new RequestFailedException("PATCH request failed for URL: " + url, e);
        }
    }
}
