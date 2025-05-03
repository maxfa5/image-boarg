package org.image_board.service;

import org.image_board.DTO.LoginRequestDTO;
import org.image_board.Model.User;
import org.image_board.repository.UserRepository;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class UserService {
    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
//        this.passwordEncoder = passwordEncoder;
    }

    public boolean checkUser(LoginRequestDTO registrationDto) throws Exception {
        User user = userRepository.findByUsername(registrationDto.getUsername())
                .orElseThrow(() -> new RuntimeException("Пользователь не найден"));

        if (!registrationDto.getPassword().equals(user.getPassword())){
            throw new RuntimeException("Неверный пароль");
        }

        return true;
    }
}