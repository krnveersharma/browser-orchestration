package com.webtest.webtest.entity;

import jakarta.persistence.ElementCollection;
import jakarta.persistence.Entity;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.Table;

import java.util.List;

@Entity
@Table(name = "sessions")
public class Session {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private long id;
    private String url;
    
    @ElementCollection
    private List<Instruction> instructions;

    @ManyToOne
    @JoinColumn(name = "user_id")
    private User user;
    private String browser;
    private String status;

    public long getId() {
        return id;
    }

    public String getUrl() {
        return url;
    }

    public List<Instruction> getInstructions() {
        return instructions;
    }

    public User getUser() {
        return user;
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

    public void setUser(User user) {
        this.user = user;
    }

    public void setBrowser(String browser) {
        this.browser = browser;
    }

    public void setStatus(String status) {
        this.status = status;
    }

}
