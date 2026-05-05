package com.threatvisualizer.backendjava.controller;

import com.threatvisualizer.backendjava.service.DataProcessingService;
import com.threatvisualizer.backendjava.service.ShieldStateService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api")
public class VulnerableController {

    @Autowired
    private DataProcessingService dataProcessingService;

    @Autowired
    private ShieldStateService shieldStateService;

    @PostMapping("/purge")
    public ResponseEntity<String> purge() {
        try {
            dataProcessingService.executeAdminPurge();
            return ResponseEntity.ok("Table purged successfully.");
        } catch (SecurityException e) {
            return ResponseEntity.status(403).body("Security Block: " + e.getMessage());
        } catch (Exception e) {
            return ResponseEntity.status(500).body("Error: " + e.getMessage());
        }
    }

    @PostMapping("/insert")
    public ResponseEntity<String> insert(@RequestParam String username, @RequestParam String data) {
        try {
            dataProcessingService.insertUserData(username, data);
            return ResponseEntity.ok("Data inserted successfully.");
        } catch (Exception e) {
            return ResponseEntity.status(500).body("Error: " + e.getMessage());
        }
    }

    @PostMapping("/toggle-shield")
    public ResponseEntity<String> toggleShield() {
        shieldStateService.toggleShield();
        return ResponseEntity.ok("Shield is now " + (shieldStateService.isActive() ? "ACTIVE" : "INACTIVE"));
    }
}
