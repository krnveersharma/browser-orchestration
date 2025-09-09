package com.webtest.webtest.allocateSession;

public interface SessionAllocationInterface {
    public void allocate(Long sessionId, String browser,String instructions, String url);
}