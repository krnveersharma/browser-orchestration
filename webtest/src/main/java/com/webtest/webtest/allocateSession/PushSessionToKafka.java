package com.webtest.webtest.allocateSession;

import com.webtest.webtest.entity.Instruction;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.webtest.webtest.kafka.MessageProducer.MessageProducer;

import java.util.List;

@Service
public class PushSessionToKafka implements SessionAllocationInterface {

    private final MessageProducer messageProducer;
    private final ObjectMapper objectMapper;

    public PushSessionToKafka(MessageProducer messageProducer, ObjectMapper objectMapper) {
        this.messageProducer = messageProducer;
        this.objectMapper = objectMapper;
    }

    @Override
    public void allocate(Long sessionId, String browser, List<Instruction> instructions, String url) {
        try {
            Session session = new Session(sessionId, browser,instructions,url);

            String json = objectMapper.writeValueAsString(session);

            messageProducer.sendMessage("session-topic", json);

            System.out.println("Sent session to Kafka: " + json);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    //making kafka message body as json
    public static class Session {

        private Long sessionId;
        private String browser;
        private List<Instruction>  instructions;
        private String url;

        public Session(Long sessionId, String browser, List<Instruction> instructions,String url) {
            this.sessionId = sessionId;
            this.browser = browser;
            this.instructions=instructions;
            this.url=url;
        }

        public Long getSessionId() {
            return sessionId;
        }

        public void setSessionId(Long sessionId) {
            this.sessionId = sessionId;
        }

        public String getBrowser() {
            return browser;
        }

        public void setBrowser(String browser) {
            this.browser = browser;
        }

        public List<Instruction>  getInstructions(){
            return instructions; 
        }

        public String getUrl(){
            return url;
        }
    }
}
