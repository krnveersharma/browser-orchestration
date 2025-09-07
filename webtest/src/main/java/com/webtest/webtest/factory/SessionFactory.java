package com.webtest.webtest.factory;

import org.springframework.stereotype.Component;

import com.webtest.webtest.dtos.SessionPostRequest;
import com.webtest.webtest.entity.Session;
import com.webtest.webtest.entity.User;

@Component
public class SessionFactory {
    public Session buildSession(SessionPostRequest request, User user) {
        Session session = new Session();
        session.setUrl(request.getUrl());
        session.setInstructions(request.getInstructions());
        session.setBrowser(request.getBrowser());
        session.setStatus("pending");
        session.setUser(user);
        return session;
    }
}
