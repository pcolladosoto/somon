import csv, datetime, openpyxl

# CSV format:
#             ts                ,      imei     , latitude, longitude, sensor,     measure
#  2025-09-08 09:02:24.077062+02,861562073825389,40.544397, -3.705246,    1  ,26.739999771118164
def processRow(row: list):
	# Drop the timezone!
	dt = datetime.datetime.strptime(row[0].split("+")[0], "%Y-%m-%d %H:%M:%S.%f")
	return [
		f"{dt.day:02}/{dt.month:02}/{dt.year:04} {dt.hour:02}:{dt.minute:02}:{dt.second:02}",
		row[1], float(row[2]), float(row[3]), int(row[4]), round(float(row[5]), 4)
	]

def createWorkbook():
	wb = openpyxl.Workbook()

	ws = wb.active
	ws.title = "temperature"
	wb.create_sheet("humidity")
	wb.create_sheet("conductivity")

	return wb

def main():
	wb = createWorkbook()

	for ws in wb:
		print(f"opening {ws.title}.csv")
		with open(f"{ws.title}.csv") as fd:
			rd = csv.reader(fd)
			for i, row in enumerate(rd):
				for j, val in enumerate(processRow(row)):
					ws.cell(row = i+1, column = j+1, value = val)

	wb.save("data.xlsx")

if __name__ == "__main__":
	main()
