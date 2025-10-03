package com.webtest.webtest.dtos;

import org.springframework.security.core.Transient;

@Transient
public class Instruction {
    private String action;
    private String value;
    private String selector;

    public String getAction() {
        return action;
    }

    public void setAction(String action) {
        this.action = action;
    }

    public String getValue() {
        return value;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public String getSelector() {
        return selector;
    }

    public void setSelector(String selector) {
        this.selector = selector;
    }
}
