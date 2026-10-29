// ทำ api
package main

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func studentRoute(app *fiber.App, db *gorm.DB) {
	//api post for create (Endpoint:/api/students)
	app.Post("/api/students", func(c fiber.Ctx) error {
		//ตัวแปรสำหรับรับข้อมูลจาก client-user
		student := new(Student)

		//นำข้อมูลจาก client-user มาเก็บในตัวแปร student พร้อมตรวจสอบ
		if error := c.Bind().Body(student); error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		//นำข้อมูลบันทึกลงตารางใน DB พร้อมตรวจสอบ
		if error := createStudent(db, student); error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed sever error",
			})
		}

		//กรณีบันทึกข้อมูลสำเร็จ ส่งกลับไปยัง client-user
		return c.Status(201).JSON(fiber.Map{
			"message": "Student created successfully",
			"data":    student,
		})
	})

	//API get for read (Endpoint:/api/students)
	app.Get("/api/students", func(c fiber.Ctx) error {
		//ตัวแปรสำหรับเก็บข้อมูลที่จะส่งกลับไปยัง client-user
		var students []Student

		//---------------------------------------------
		//สร้างตัวแปรสำหรับเก็บข้อมูลที่เป็นเงื่อนไขที่ส่งมาแบ query string
		StuName := c.Query("StuName")
		StuLevel := c.Query("StuLevel")

		var error error

		if StuName != "" {
			//ดึงข้อมูลเงื่อนไข StuName
			students, error = getStudentByName(db, StuName)
		} else if StuLevel != "" {
			//ดึงข้อมูลเงื่อนไข StuLevel
			StuLevelInt, err := strconv.Atoi(StuLevel)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{
					"message": "Invalid StuLevel parameter",
				})
			}
			students, error = getStudentByLevel(db, StuLevelInt)
		} else {
			//ดึงข้อมูลจากตาราง เก็บในตัวแปร
			students, error = getAllStudent(db)
		}

		//---------------------------------------------

		//ตรวจสอบว่ามีข้อมูลหรือไม่
		if errors.Is(error, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"message": "Student not found",
			})
		}

		//ตรวจสอบว่ามีการทำงานผิดพลาดหรือไม่
		if error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to retrieve students",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"message": "Students retrieved successfully",
			"data":    students,
		})
	})

	//API get for read by id (Endpoint:/api/students/:id)
	//เงื่อนไขส่งมาแบบ path parameter ส่งต่อท้าย endpoint
	app.Get("/api/students/:id", func(c fiber.Ctx) error {
		//สร้างตัวเก็บ PK คือ ID ที่ส่งมา
		id := c.Params("id")
		//แปลง id เป็นตตัวเลข
		idInt, error := strconv.Atoi(id)
		if error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request parameter",
			})
		}

		//ดึงข้อมูลเก็บในตัวแปร
		student, error := getStudentByID(db, idInt)

		//ตรวจสอบว่ามีข้อมูลหรือไม่
		if errors.Is(error, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{
				"message": "Student not found",
			})
		}

		//มีการทำงานผิดพลาดหรือไม่
		if error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to retrieve student",
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"message": "Student retrieved successfully",
			"data":    student,
		})
	})

	//API put for update (Endpoint:/api/students/:id)
	app.Put("/api/students/:id", func(c fiber.Ctx) error {
		//สร้างตัวเก็บ PK คือ ID ที่ส่งมา
		id := c.Params("id")

		//แปลง id เป็นตตัวเลข
		idInt, error := strconv.Atoi(id)
		if error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request parameter",
			})
		}

		//ตัวแปรสำหรับรับข้อมูลจาก client-user
		student := new(Student)

		//นำข้อมูลที่ส่งมาเก็บในตัวแปร พร้อมตรวจสอบ
		if error := c.Bind().Body(student); error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		//นำข้อมูลไปแก้ไขในตาราง พร้อมตรวจสอบ
		if error := updateStudent(db, idInt, student); error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to update student",
			})
		}

		//กรณีแก้ไขข้อมูลสำเร็จ
		return c.Status(200).JSON(fiber.Map{
			"message": "Student updated successfully",
			"data":    student,
		})
	})

	//API delete for delete (Endpoint:/api/students/:id)
	app.Delete("/api/students/:id", func(c fiber.Ctx) error {
		//สร้างตัวเก็บ PK คือ ID ที่ส่งมา
		id := c.Params("id")

		//แปลง id เป็นตตัวเลข
		idInt, error := strconv.Atoi(id)
		if error != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Invalid request parameter",
			})
		}

		//นำข้อมูลไปลบในตาราง พร้อมตรวจสอบ
		if error := deleteStudent(db, idInt); error != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "Failed to delete student",
			})
		}

		//กรณีลบข้อมูลสำเร็จ
		return c.Status(200).JSON(fiber.Map{
			"message": "Student deleted successfully",
		})
	})

}
