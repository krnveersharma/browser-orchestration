package com.webtest.webtest.allocateSession;

import com.webtest.webtest.entity.Instruction;

import java.util.List;

public interface SessionAllocationInterface {
    public void allocate(Long sessionId, String browser, List<Instruction> instructions, String url);
}