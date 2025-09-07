package com.webtest.webtest.repository;

import org.springframework.data.jpa.repository.JpaRepository;

import com.webtest.webtest.entity.Session;

public interface  SessionRepository extends JpaRepository<Session, Long> {
}
