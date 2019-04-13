//
//  Permission_HubTests.swift
//  Permission HubTests
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import XCTest
@testable import Permission_Hub

class Permission_HubTests: XCTestCase {

    override func setUp() {
        // Put setup code here. This method is called before the invocation of each test method in the class.
    }

    override func tearDown() {
        // Put teardown code here. This method is called after the invocation of each test method in the class.
    }

    func testExample() {
        // This is an example of a functional test case.
        // Use XCTAssert and related functions to verify your tests produce the correct results.

        testDataParsing()
    }

    func testDataParsing() {

        let jsonString = "{\"transactionID\":\"b80db272b05b9ad007c6833dac68b95ca907594946b2da1929d1f8f95d973b5c\",\"item\":[{\"Item\":\"Personal details\",\"Fields\":[\"name\",\"surname\",\"birth_date\",\"email\",\"BSN\"]},{\"Item\":\"Passport\",\"Fields\":[\"number\",\"expiration\",\"country\"]},{\"Item\":\"Banking details\",\"Fields\":[\"IBAN\",\"bank\",\"name\"]}],\"title\":\"Provide permission for completing\",\"description\":\"T-mobile monthly plan(unlimited data), 65 euro, iPhone XR 256GB\",\"verification\":[\"digid.nl\",\"planet-blockchain\",\"kvk\"],\"date\":\"2019-04-13T15:51:57+02:00\",\"requesterName\":\"John Smith\",\"RequesterPublicKey\":\"69093eef7426963f2ef0f68fb73e355b7898ddb04a4fad769a96b41ffc824c1c\",\"analysis\":[\"personal data is GDPR protected data\",\"banking details is sensitive data\"]}"

        let jsonData = jsonString.data(using: .utf8)!
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601

        do {
            let transaction = try decoder.decode(TransactionNotification.self, from: jsonData)
            print(transaction)
            XCTAssert(true, "JSON parsing succeeded")
        } catch {
            XCTAssert(false, "JSON parsing failed")
        }
    }

    func testPerformanceExample() {
        // This is an example of a performance test case.
        self.measure {
            // Put the code you want to measure the time of here.
        }
    }

}
