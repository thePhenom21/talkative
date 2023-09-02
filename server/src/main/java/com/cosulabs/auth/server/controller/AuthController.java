package com.cosulabs.auth.server.controller;

import com.cosulabs.auth.server.models.UserEntity;
import com.cosulabs.auth.server.repository.UserEntityRepository;
import com.cosulabs.auth.server.services.JwtService;
import org.apache.coyote.Response;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.bind.annotation.*;

@RestController
public class AuthController {

    private UserEntityRepository userEntityRepository;

    private JwtService jwtService;

    private PasswordEncoder passwordEncoder;

    AuthController(UserEntityRepository userEntityRepository, JwtService jwtService){
        this.userEntityRepository = userEntityRepository;
        this.jwtService = jwtService;
    }



    @PostMapping("/register/{email}/{username}/{password}")
    public ResponseEntity<String> register(@PathVariable String username, @PathVariable String password, @PathVariable String email){
        try{
            UserEntity user = new UserEntity();
            user.setUsername(username);
            user.setPassword(password);
            user.setEmail(email);

            String token = jwtService.generateToken(user);
            user.setToken(token);

            userEntityRepository.save(user);
            return ResponseEntity.ok("User created with username "+username);
        }catch (Exception e){

        }
        return ResponseEntity.badRequest().build();
    }

    @PostMapping("/login/{username}/{password}")
    public ResponseEntity<String> login(@PathVariable String username, @PathVariable String password){
        try{
            UserEntity user = userEntityRepository.findByUsername(username);
            if(passwordEncoder.encode(password).equals(user.getPassword())){
                return ResponseEntity.ok().body(user.getToken());
            }
        } catch (Exception e){

        }
        return ResponseEntity.badRequest().build();
    }

    @PostMapping("/auth")
    public ResponseEntity<String> auth(@RequestBody String body){
        try{
        UserEntity usr = userEntityRepository.findUserEntityByToken(body);
        if(jwtService.isTokenValid(body,usr)){
            return ResponseEntity.ok("User logged in: "+usr.getUsername());
        }}
        catch (Exception e){

        }
        return ResponseEntity.badRequest().body("Not authorized");
    }



}
