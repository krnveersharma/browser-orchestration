package com.webtest.webtest.entity;

import jakarta.persistence.Embeddable;

@Embeddable
public class Instruction {
    private String action;
    private String value;
    private String selector;

    // Default constructor required by JPA
    public Instruction() {}

    public Instruction(String action, String value, String selector) {
        this.action = action;
        this.value = value;
        this.selector = selector;
    }

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
