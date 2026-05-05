package com.threatvisualizer.backendjava;

import org.springframework.stereotype.Service;

@Service
public class ShieldStateService {
    private boolean active = false;

    public boolean isActive() {
        return active;
    }

    public void setActive(boolean active) {
        this.active = active;
    }

    public void toggleShield() {
        this.active = !this.active;
    }
}
