use bevy::prelude::*;
use bevy::input::ButtonInput;
use bevy::input::keyboard::KeyCode;
use bevy_rapier3d::prelude::*;

// Strongly typed marker for the player character
#[derive(Component)]
struct Player;

fn setup_graphics(mut commands: Commands) {
    // Add a camera so we can see the debug-render.
    commands.spawn((
        Camera3d::default(),
        Transform::from_xyz(-3.0, 3.0, 10.0).looking_at(Vec3::ZERO, Vec3::Y),
    ));
}

fn setup_physics(mut commands: Commands) {
    // Create the ground (flat plane)
    commands
        .spawn(Collider::cuboid(100.0, 0.1, 100.0))
        .insert(Transform::from_xyz(0.0, -2.0, 0.0));

    // Create the kinematic character controller (capsule shape) with Player marker
    commands
        .spawn((
            RigidBody::KinematicPositionBased,
            Collider::capsule_y(1.0, 0.5),
            Transform::from_xyz(0.0, 0.0, 0.0),
            KinematicCharacterController::default(),
            Player,
        ));
}

fn character_movement(
    keyboard_input: Res<ButtonInput<KeyCode>>,
    time: Res<Time>,
    mut controllers: Query<&mut KinematicCharacterController, With<Player>>,
) {
    let mut direction: Vec3 = Vec3::ZERO;
    if keyboard_input.pressed(KeyCode::KeyW) {
        direction.z -= 1.0;
    }
    if keyboard_input.pressed(KeyCode::KeyS) {
        direction.z += 1.0;
    }
    if keyboard_input.pressed(KeyCode::KeyA) {
        direction.x -= 1.0;
    }
    if keyboard_input.pressed(KeyCode::KeyD) {
        direction.x += 1.0;
    }
    if direction.length_squared() > 0.0 {
        direction = direction.normalize();
    }
    let speed: f32 = 5.0;
    let movement: Vec3 = direction * speed * time.delta_secs();
    for mut controller in controllers.iter_mut() {
        controller.translation = Some(movement);
    }
}

fn main() {
    App::new()
        .add_plugins(DefaultPlugins)
        .add_plugins(RapierPhysicsPlugin::<NoUserData>::default())
        .add_plugins(RapierDebugRenderPlugin::default())
        .add_systems(Startup, setup_graphics)
        .add_systems(Startup, setup_physics)
        .add_systems(Update, character_movement)
        .run();
}