package com.webtest.webtest.dtos;

public class SessionPostRequest {

    private String url;
    private String instructions;
    private String browser;
    private String status;

    public String getUrl() {
        return url;
    }

    public String getInstructions() {
        return instructions;
    }

    public String getBrowser() {
        return browser;
    }

    public String getStatus() {
        return status;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public void setInstructions(String instructions) {
        this.instructions = instructions;
    }

    public void setBrowser(String browser) {
        this.browser = browser;
    }

    public void setStatus(String status) {
        this.status = status;
    }
}
