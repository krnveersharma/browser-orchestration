package com.webtest.webtest.service;

import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

import com.webtest.webtest.entity.User;
import com.webtest.webtest.repository.UserRepository;
import com.webtest.webtest.util.JwtUtil;


@Service
public class AuthService {
    private final UserRepository userRepository;
    private final BCryptPasswordEncoder passwordEncoder;
    private final JwtUtil jwtUtil;

    public AuthService(UserRepository userRepository) {
        this.userRepository = userRepository;
        this.passwordEncoder = new BCryptPasswordEncoder();
        this.jwtUtil=new JwtUtil();
    }

    public User signup(String name, String email, String password){
        String hashedPassword=passwordEncoder.encode(password);
        User user=new User();
        user.setEmail(email);
        user.setName(name);
        user.setPassword(hashedPassword);

        return userRepository.save(user);
    }

   public String login(String email, String password) throws Exception {
        User user = userRepository.findByEmail(email)
                .orElseThrow(() -> new Exception("User not found"));

        if (!passwordEncoder.matches(password, user.getPassword())) {
            throw new Exception("Invalid credentials");
        }

        return jwtUtil.generateToken(user.getEmail());
    }

    public User getUserByEmail(String email){
        return userRepository.findByEmail(email).orElse(null);
    }
}
