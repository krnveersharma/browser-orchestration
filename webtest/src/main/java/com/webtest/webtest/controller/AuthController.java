package com.webtest.webtest.controller;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.webtest.webtest.dtos.LoginRequest;
import com.webtest.webtest.dtos.SignupRequest;
import com.webtest.webtest.entity.User;
import com.webtest.webtest.service.AuthService;

@RestController
@RequestMapping("/auth")
public class AuthController {

    private final AuthService authService;

    public AuthController(AuthService authService) {
        this.authService = authService;
    }

    @PostMapping("/login")
    public ResponseEntity<?> login(@RequestBody LoginRequest request) {
        try {
            String token = authService.login(request.getEmail(), request.getPassword());
            return ResponseEntity.ok(token);
        } catch (Exception e) {
            String message = e.getMessage();
            HttpStatus status;
            if (message.contains("User not found")) {
                status = HttpStatus.NOT_FOUND;
            } else if (message.contains("Invalid credentials")) {
                status = HttpStatus.UNAUTHORIZED;
            } else {
                status = HttpStatus.BAD_REQUEST;
            }

            return ResponseEntity.status(status).body(message);
        }
    }

    @PostMapping("/sign-up")
    public ResponseEntity<?> signup(@RequestBody SignupRequest request) {
        try {
            User user = authService.signup(request.getName(), request.getEmail(), request.getPassword());
            return ResponseEntity.status(HttpStatus.CREATED).body(user);
        } catch (Exception e) {
            String message = e.getMessage();
            HttpStatus status;

            if (message.contains("Email already exists")) {
                status = HttpStatus.CONFLICT;
             }else {
                status = HttpStatus.BAD_REQUEST;
            }

            return ResponseEntity.status(status).body(message);
        }
    }
}
