package com.webtest.webtest.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.webtest.webtest.dtos.SessionPostRequest;
import com.webtest.webtest.entity.Session;
import com.webtest.webtest.service.SessionService;

@RestController
@RequestMapping("/session")
public class SessionController {
    private final SessionService sessionService;

    public SessionController(SessionService sessionService){
        this.sessionService=sessionService;
    }
    @PostMapping("/create-session")
    public ResponseEntity<?> CreateSession(@RequestBody SessionPostRequest sessionRequest,@RequestHeader("Authorization") String token){
        Session session=sessionService.createSession(sessionRequest, token);
        return  ResponseEntity.status(400).body(session);
    }
}
