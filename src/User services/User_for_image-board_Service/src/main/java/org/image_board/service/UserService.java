package org.image_board.service;

import jakarta.transaction.Transactional;
import org.image_board.DTO.UserRegistrationDto;
import org.image_board.Model.User;
import org.image_board.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.List;
import java.util.Optional;
import java.util.Set;

@Service
public class UserService {
    private final UserRepository userRepository;

    public UserService(UserRepository userRepository) {
        this.userRepository = userRepository;
//        this.passwordEncoder = passwordEncoder;
    }

    public User registerNewUser(UserRegistrationDto registrationDto) throws Exception {
        if (userRepository.existsByUsername(registrationDto.getUsername())) {
            throw new Exception("Пользователь с таким именем уже существует");
        }

        User user = new User(
                registrationDto.getUsername(),
                registrationDto.getPassword()
        );

        return userRepository.save(user);
    }

    public List<User> getAllUsers() {
        return userRepository.findAll(); // Use Spring Data JPA to fetch all users
    }

}