package org.image_board.controller;


import org.image_board.DTO.LoginRequestDTO;
import org.image_board.Utils.JwtUtils;
import org.image_board.service.UserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

import java.util.List;

@RestController
public class AuthController {
//    @PostMapping("/auth/login")
//    public ResponseEntity<String> login(@RequestBody LoginRequestDTO request) {
//        System.out.println("Login request: " + request);
//        // Проверка логина/пароля (можно использовать UserDetailsService)
//        String token = JwtUtils.generateToken(request.getUsername());
//        return ResponseEntity.ok(token);
//    }
    @Autowired
    public AuthController(UserService userService) {
        this.userService = userService;
    }
    private UserService userService;

    @PostMapping("auth/login")
    public ResponseEntity<?> login(@RequestBody LoginRequestDTO request) {
        try {
            if (userService.checkUser(request)) {
                String token = JwtUtils.generateToken(request.getUsername());
                return ResponseEntity.ok(token);
            }
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();
        } catch (Exception e) {
                return ResponseEntity.status(HttpStatus.UNAUTHORIZED)
                    .body(e.getMessage());
        }
    }
}