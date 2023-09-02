package com.cosulabs.auth.server.repository;

import com.cosulabs.auth.server.models.UserEntity;
import org.springframework.data.jpa.repository.JpaRepository;

public interface UserEntityRepository extends JpaRepository<UserEntity,Long> {

    public UserEntity findByEmail(String email);

    public UserEntity findByUsername(String username);

    public UserEntity findUserEntityByToken(String token);

}
