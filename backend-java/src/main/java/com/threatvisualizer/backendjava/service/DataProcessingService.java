package com.threatvisualizer.backendjava.service;

import com.threatvisualizer.backendjava.model.UserData;
import com.threatvisualizer.backendjava.repository.SecureDataRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.UUID;

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
    public void executeSqliAttack() {
        if (shieldStateService.isActive()) {
            throw new SecurityException("SQL Injection attack is blocked by the shield.");
        } else {
            String randomPwned = "PWNED_" + (int)(Math.random() * 9000 + 1000);
            String sql = "UPDATE user_data SET username = '" + randomPwned + "'";
            jdbcTemplate.execute(sql);
        }
    }

    @Transactional(readOnly = true)
    public List<UserData> getAllData() {
        return secureDataRepository.findAll();
    }

    @Transactional
    public UserData seedData() {
        UserData userData = new UserData();
        userData.setUsername("TestUser_" + UUID.randomUUID().toString().substring(0, 8));
        userData.setData("System secure");
        return secureDataRepository.save(userData);
    }

    @Transactional
    public void clearAllData() {
        secureDataRepository.deleteAll();
    }
}
