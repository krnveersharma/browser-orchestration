package com.webtest.webtest.util;

import java.security.Key;
import java.util.Date;

import org.springframework.stereotype.Component;

import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.SignatureAlgorithm;
import io.jsonwebtoken.security.Keys;

@Component
public class JwtUtil {
    private final String SECRET_KEY = "my_super_secret_key_change_this_to_env_variable_123456";
    private final long EXPIRATION = 1000 * 60 * 60 * 10; 
    private final Key key = Keys.hmacShaKeyFor(SECRET_KEY.getBytes());

    public String generateToken(String email) {
        return Jwts.builder()
                .setSubject(email)
                .setIssuedAt(new Date())
                .setExpiration(new Date(System.currentTimeMillis() + EXPIRATION))
                .signWith(key, SignatureAlgorithm.HS256)
                .compact();
    }

    public String GetEmail(String authHeader) throws Exception{
        try {
            String token = authHeader.replace("Bearer ", "");

            Claims claims = Jwts.parserBuilder()
                    .setSigningKey(key)
                    .build()
                    .parseClaimsJws(token)
                    .getBody();
            System.out.println("email is: "+ claims.getSubject());
            return claims.getSubject();
        } catch (Exception e) {
            e.printStackTrace();
            throw e;
        }
    }
}
