package com.messenger.android.ui.auth

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AuthScreen(
    viewModel: AuthViewModel = hiltViewModel(),
    onNavigateToChat: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()
    
    // Навигация при успешной авторизации
    LaunchedEffect(uiState.isAuthenticated) {
        if (uiState.isAuthenticated) {
            onNavigateToChat()
        }
    }
    
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp)
            .verticalScroll(rememberScrollState()),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Spacer(modifier = Modifier.height(64.dp))
        
        // Заголовок
        Text(
            text = "Messenger",
            fontSize = 32.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.primary
        )
        
        Text(
            text = "Современный мессенджер для общения",
            fontSize = 16.sp,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            textAlign = TextAlign.Center,
            modifier = Modifier.padding(top = 8.dp, bottom = 48.dp)
        )
        
        // Переключатель между входом и регистрацией
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            TextButton(
                onClick = { viewModel.setAuthMode(AuthMode.LOGIN) },
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = "Вход",
                    color = if (uiState.authMode == AuthMode.LOGIN) {
                        MaterialTheme.colorScheme.primary
                    } else {
                        MaterialTheme.colorScheme.onSurfaceVariant
                    }
                )
            }
            
            TextButton(
                onClick = { viewModel.setAuthMode(AuthMode.REGISTER) },
                modifier = Modifier.weight(1f)
            ) {
                Text(
                    text = "Регистрация",
                    color = if (uiState.authMode == AuthMode.REGISTER) {
                        MaterialTheme.colorScheme.primary
                    } else {
                        MaterialTheme.colorScheme.onSurfaceVariant
                    }
                )
            }
        }
        
        Spacer(modifier = Modifier.height(24.dp))
        
        // Форма входа/регистрации
        Card(
            modifier = Modifier.fillMaxWidth(),
            elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
        ) {
            Column(
                modifier = Modifier.padding(16.dp),
                verticalArrangement = Arrangement.spacedBy(16.dp)
            ) {
                // Поле имени пользователя
                OutlinedTextField(
                    value = uiState.username,
                    onValueChange = viewModel::setUsername,
                    label = { Text("Имя пользователя") },
                    modifier = Modifier.fillMaxWidth(),
                    isError = uiState.usernameError != null,
                    supportingText = uiState.usernameError?.let { { Text(it) } }
                )
                
                // Поле email (только для регистрации)
                if (uiState.authMode == AuthMode.REGISTER) {
                    OutlinedTextField(
                        value = uiState.email,
                        onValueChange = viewModel::setEmail,
                        label = { Text("Email") },
                        modifier = Modifier.fillMaxWidth(),
                        isError = uiState.emailError != null,
                        supportingText = uiState.emailError?.let { { Text(it) } }
                    )
                }
                
                // Поле пароля
                OutlinedTextField(
                    value = uiState.password,
                    onValueChange = viewModel::setPassword,
                    label = { Text("Пароль") },
                    modifier = Modifier.fillMaxWidth(),
                    visualTransformation = PasswordVisualTransformation(),
                    isError = uiState.passwordError != null,
                    supportingText = uiState.passwordError?.let { { Text(it) } }
                )
                
                // Дополнительные поля для регистрации
                if (uiState.authMode == AuthMode.REGISTER) {
                    OutlinedTextField(
                        value = uiState.firstName,
                        onValueChange = viewModel::setFirstName,
                        label = { Text("Имя") },
                        modifier = Modifier.fillMaxWidth()
                    )
                    
                    OutlinedTextField(
                        value = uiState.lastName,
                        onValueChange = viewModel::setLastName,
                        label = { Text("Фамилия") },
                        modifier = Modifier.fillMaxWidth()
                    )
                }
                
                // Кнопка отправки
                Button(
                    onClick = {
                        if (uiState.authMode == AuthMode.LOGIN) {
                            viewModel.login()
                        } else {
                            viewModel.register()
                        }
                    },
                    modifier = Modifier.fillMaxWidth(),
                    enabled = !uiState.isLoading
                ) {
                    if (uiState.isLoading) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(16.dp),
                            color = MaterialTheme.colorScheme.onPrimary
                        )
                    } else {
                        Text(
                            text = if (uiState.authMode == AuthMode.LOGIN) "Войти" else "Зарегистрироваться"
                        )
                    }
                }
            }
        }
        
        // Показ ошибок
        if (uiState.error != null) {
            Spacer(modifier = Modifier.height(16.dp))
            Card(
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.errorContainer
                )
            ) {
                Text(
                    text = uiState.error,
                    color = MaterialTheme.colorScheme.onErrorContainer,
                    modifier = Modifier.padding(16.dp)
                )
            }
        }
        
        Spacer(modifier = Modifier.height(32.dp))
    }
}