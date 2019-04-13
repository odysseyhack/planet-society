//
//  TransactionOverviewViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionOverviewViewController: PHTableViewController {

    // MARK: - Private properties

    private var transaction: TransactionNotification

    // MARK: - Life cycle

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        var items = [PHTableViewViewCellType]()
        items.append(.notification(
            type: .verification,
            text: "This company is verified"))
        items.append(.notification(
            type: .warning,
            text: "Permission warning!"))
        items.append(.description(
            date: transaction.date,
            title: transaction.title,
            description: transaction.description))

        transaction.items.forEach {
            items.append(.transactionItem(item: $0))
        }

        super.init(
            title: transaction.requesterName,
            items: items)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()

        // set selection delegate to self
        delegate = self
    }
}

extension TransactionOverviewViewController: PHTableViewControllerDelegate {

    func didSelect(item: PHTableViewViewCellType) {

        switch item {
        case .notification(let type, let text):

            let items: [PHTableViewViewCellType]
            switch type {
            case .warning:
                items = transaction.analysis.map { .description(
                    date: Date(),
                    title: "",
                    description: $0) }
                let viewController = PHTableViewController(
                    title: text,
                    items: items)

                let notification = Notification(
                    name: Notification.Name("show warning"),
                    object: viewController,
                    userInfo: nil)
                NotificationCenter.default.post(notification)

            case .verification:
                 items = transaction.verification.map { .description(
                    date: Date(),
                    title: "",
                    description: $0) }

                 let viewController = PHTableViewController(
                    title: text,
                    items: items)

                 let notification = Notification(
                    name: Notification.Name("show verification"),
                    object: viewController,
                    userInfo: nil)
                 NotificationCenter.default.post(notification)
            }

        case .transactionItem(let item):
            if let index = self.transaction.items.index(of: item) {
                transaction.items[index].isAccepted = true
            }
//            validateSelection()

        default:
            break
        }
    }

    func didDeselect(item: PHTableViewViewCellType) {

        switch item {
        case .transactionItem(let item):
            if let index = self.transaction.items.index(of: item) {
                transaction.items[index].isAccepted = false
            }
//            validateSelection()

        default:
            break
        }
    }
}
