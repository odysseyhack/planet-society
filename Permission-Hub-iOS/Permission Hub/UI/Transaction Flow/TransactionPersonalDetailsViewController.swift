//
//  TransactionPersonalDetailsViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionPersonalDetailsViewController: PHTableViewController {

    // MARK: - Private properties

    private var transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(title: "", items: [
            .notification(
                type: .verification,
                text: "This company is verified"),
            .description(
                date: transaction.date,
                title: "Personal details",
                description: "Please fill out your personal details."),
            .plugin(
                image: UIImage(named: "digid_button"),
                text: "Use the external DigiD plug-in to fill in your personal information (optional)."),
            .form(placeholder: "First name", text: nil),
            .form(placeholder: "Last name", text: nil),
            .form(placeholder: "Date of birth", text: nil),
            .form(placeholder: "Address", text: nil),
            .form(placeholder: "Email", text: nil),
            .form(placeholder: "BSN number", text: nil)
        ])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        tableView.allowsSelection = false
    }
}
