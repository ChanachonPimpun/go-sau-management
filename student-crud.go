// ทำงานกับ DB-CRUD
package main

import "gorm.io/gorm"

//ฟังก์ชันเพิ่มข้อมูล
func createStudent(db *gorm.DB, student *Student) error {
	return db.Create(student).Error
}

//ฟังก์ชันดึงข้อมูลจากตาราง
func getAllStudent(db *gorm.DB) ([]Student, error) {
	//สร้างตัวแปรสำหรับเก็บข้อมูลจากตาราง
	var students []Student
	result := db.Find(&students)
	return students, result.Error
}

//ฟังก์ชันดึงข้อมูลจากตารางตาม เงื่อนไข PK คือ ID
func getStudentByID(db *gorm.DB, id int) (*Student, error) {
	//สร้างตัวแปรสำหรับเก็บข้อมูลจากตาราง
	var students Student

	result := db.First(&students, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &students, nil
}

//ฟังก์ชันแก้ไขข้อมูลจากตารางตาม เงื่อนไขไม่ใช่ PK เช่น ชื่อ ตรงกับที่ส่งมาบางส่วน
func getStudentByName(db *gorm.DB, StuName string) ([]Student, error) {
	//สร้างตัวแปรสำหรับเก็บข้อมูลจากตาราง
	var students []Student

	result := db.Where("Stu_name LIKE ?", "%"+StuName+"%").Find(&students)

	return students, result.Error
}

//ฟังก์ชันดึงข้อมูลจากตารางตาม เงื่อนไขชั้นปี
func getStudentByLevel(db *gorm.DB, StuLevel int) ([]Student, error) {
	//สร้างตัวแปรสำหรับเก็บข้อมูลจากตาราง
	var students []Student

	result := db.Where("Stu_level = ?", StuLevel).Find(&students)

	return students, result.Error
}

//ฟังก์ชันแก้ไขข้อมูลจากตารางตาม เงื่อนไข PK คือ ID
func updateStudent(db *gorm.DB, id int, student *Student) error {
	return db.Model(&Student{}).Where("id = ?", id).Updates(student).Error
}

//ฟังก์ชันลบข้อมูลจากตารางตาม เงื่อนไข PK คือ ID
func deleteStudent(db *gorm.DB, id int) error {
	return db.Delete(&Student{}, id).Error
}
