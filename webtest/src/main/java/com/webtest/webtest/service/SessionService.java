package com.webtest.webtest.service;

import org.springframework.stereotype.Service;

import com.webtest.webtest.allocateSession.PushSessionToKafka;
import com.webtest.webtest.allocateSession.SessionAllocationInterface;
import com.webtest.webtest.dtos.SessionPostRequest;
import com.webtest.webtest.entity.Session;
import com.webtest.webtest.entity.User;
import com.webtest.webtest.exception.UserNotFoundException;
import com.webtest.webtest.factory.SessionFactory;
import com.webtest.webtest.repository.SessionRepository;
import com.webtest.webtest.util.JwtUtil;

@Service
public class SessionService {

    private final SessionRepository sessionRepository;
    private final JwtUtil jwtUtil;
    private final AuthService authService;
    private final SessionFactory sessionFactory;
    private final SessionAllocationInterface allocationStrategy; 

    public SessionService(SessionRepository sessionRepository, AuthService authService, SessionFactory sessionFactory,PushSessionToKafka allocationStrategy) {
        this.sessionRepository = sessionRepository;
        this.jwtUtil = new JwtUtil();
        this.authService = authService;
        this.sessionFactory = sessionFactory;
        this.allocationStrategy=allocationStrategy;
    }

    public Session createSession(SessionPostRequest request, String token) {
        String email;
        try {
            email = jwtUtil.GetEmail(token);
        } catch (Exception e) {
            throw new RuntimeException("Failed to extract email from token", e);
        }
        User user = authService.getUserByEmail(email);
        if (user == null) {
            throw new UserNotFoundException(email);
        }

        Session session = sessionFactory.buildSession(request, user);
        allocationStrategy.allocate(session.getId(), session.getBrowser());
        return sessionRepository.save(session);
    }
}
