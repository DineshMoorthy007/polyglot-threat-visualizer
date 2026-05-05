package com.threatvisualizer.backendjava.service;

import com.threatvisualizer.backendjava.model.UserData;
import com.threatvisualizer.backendjava.repository.SecureDataRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class DataProcessingService {

    @Autowired
    private ShieldStateService shieldStateService;

    @Autowired
    private SecureDataRepository secureDataRepository;

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Transactional
    public void executeAdminPurge() {
        if (shieldStateService.isActive()) {
            throw new SecurityException("Admin purge is blocked by the shield.");
        } else {
            jdbcTemplate.execute("TRUNCATE TABLE user_data");
        }
    }

    @Transactional
    public void insertUserData(String username, String userInput) {
        if (shieldStateService.isActive()) {
            UserData userData = new UserData();
            userData.setUsername(username);
            userData.setData(userInput);
            secureDataRepository.save(userData);
        } else {
            String sql = "INSERT INTO user_data (username, data) VALUES ('" + username + "', '" + userInput + "')";
            jdbcTemplate.execute(sql);
        }
    }
}
