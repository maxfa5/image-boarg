package org.image_board.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;


@RestController
@RequestMapping("/api")
public class ProtectedController {
    @GetMapping("/data")
    public ResponseEntity<String> getProtectedData() {
        // Только для аутентифицированных пользователей

        return ResponseEntity.ok("Secret data");
    }

}
