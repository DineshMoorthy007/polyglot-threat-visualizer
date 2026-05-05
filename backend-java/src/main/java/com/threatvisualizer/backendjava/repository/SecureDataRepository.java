package com.threatvisualizer.backendjava.repository;

import com.threatvisualizer.backendjava.model.UserData;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface SecureDataRepository extends JpaRepository<UserData, Long> {
}
