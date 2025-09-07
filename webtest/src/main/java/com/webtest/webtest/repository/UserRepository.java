package com.webtest.webtest.repository;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;

import com.webtest.webtest.entity.User;

public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);
}
