package com.webtest.webtest.dtos;
import com.webtest.webtest.entity.Instruction;

import java.util.List;
public class SessionPostRequest {

    private String url;
    private List<Instruction> instructions;
    private String browser;
    private String status;

    public String getUrl() {
        return url;
    }

    public List<Instruction> getInstructions() {
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

    public void setInstructions(List<Instruction> instructions) {
        this.instructions = instructions;
    }

    public void setBrowser(String browser) {
        this.browser = browser;
    }

    public void setStatus(String status) {
        this.status = status;
    }
}