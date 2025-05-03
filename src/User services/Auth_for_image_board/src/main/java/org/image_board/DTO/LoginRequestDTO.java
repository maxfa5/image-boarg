package org.image_board.DTO;

import lombok.Data;

@Data
public class LoginRequestDTO {
    private String username;
    private String password;
}